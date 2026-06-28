package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/constants"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/pkg/utils"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type VoucherOrderService interface {
	// SeckillVoucher 抢购秒杀券。
	SeckillVoucher(ctx context.Context, voucherID int64, userID int64) result.Result
	// Start 启动后台订单消费者。
	Start(ctx context.Context)
}

// 新增包级脚本对象
var seckillScript = redis.NewScript(constants.SeckillScript)

type voucherOrderService struct {
	orderRepo   repository.VoucherOrderRepository // orderRepo 负责订单表操作。
	redisClient *redis.Client                     // redisClient 后面用于 Lua 扣库存、消息队列等。
}

// NewVoucherOrderService 创建优惠券订单 Service。
// 构造 Service 时启动后台消费者
func NewVoucherOrderService(
	orderRepo repository.VoucherOrderRepository,
	redisClient *redis.Client,
) VoucherOrderService {
	return &voucherOrderService{
		orderRepo:   orderRepo,
		redisClient: redisClient,
	}
}

func (s *voucherOrderService) Start(ctx context.Context) {
	s.initSeckillStream(ctx)
	go s.handleVoucherOrder(ctx)
}

// 初始化 Stream 消费者组：
func (s *voucherOrderService) initSeckillStream(ctx context.Context) {
	err := s.redisClient.XGroupCreateMkStream(
		ctx,
		constants.SeckillStreamKey,
		constants.SeckillConsumerGroup,
		"0",
	).Err()

	//BUSYGROUP 表示消费者组已经存在，不算错误。
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		log.Printf("[seckill] create stream group failed: %v", err)
	}
}

// SeckillVoucher 是秒杀下单入口。
//
// 后面会实现：Lua 判断库存和一人一单 -> 返回订单 id -> 异步写订单。
func (s *voucherOrderService) SeckillVoucher(ctx context.Context, voucherID int64, userID int64) result.Result {
	if voucherID <= 0 {
		return result.Fail("invalid voucher id")
	}
	if userID <= 0 {
		return result.Fail("用户未登录")
	}
	orderID := utils.GlobalSnowflake.NextVal()

	code, err := seckillScript.Run(
		ctx,
		s.redisClient,
		[]string{},
		strconv.FormatInt(voucherID, 10),
		strconv.FormatInt(userID, 10),
		strconv.FormatInt(orderID, 10),
	).Int()

	if err != nil {
		log.Printf("[seckill] lua execute failed, voucherID=%d userID=%d err=%v", voucherID, userID, err)
		return result.Fail("秒杀失败，请稍后重试")
	}

	switch code {
	case 0:
		return result.OKWithData(orderID)
	case 1:
		return result.Fail("库存不足")
	case 2:
		return result.Fail("不能重复下单")
	default:
		return result.Fail("秒杀失败")
	}
}

// 后台 goroutine 消费 Stream
//
// 每轮循环开始先判断 ctx 是否取消,
// 如果取消，就退出,
// 如果已经读到消息，就处理完当前消息再 ACK。
func (s *voucherOrderService) handleVoucherOrder(ctx context.Context) {

	// 服务启动后先处理历史 pending 消息，避免之前宕机留下的消息永远不被 ACK。
	s.handlePendingList(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("[seckill] voucher order consumer stopped")
			return
		default:
		}

		streams, err := s.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    constants.SeckillConsumerGroup,
			Consumer: constants.SeckillConsumerName,
			Streams:  []string{constants.SeckillStreamKey, ">"},
			Count:    1,
			Block:    2 * time.Second,
		}).Result()

		if errors.Is(err, redis.Nil) {
			// 空闲时也可以顺手扫一下 pending。
			s.handlePendingList(ctx)
			continue
		}

		if err != nil {
			log.Printf("[seckill] read stream failed: %v", err)
			time.Sleep(500 * time.Millisecond)
			s.handlePendingList(ctx)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				s.handleMessageWithRetry(ctx, msg)
			}
		}
	}
}

// 处理单条订单消息
func (s *voucherOrderService) handleVoucherOrderMessage(ctx context.Context, msg redis.XMessage) error {
	order, err := parseVoucherOrderMessage(msg)
	if err != nil {
		return err
	}

	return s.orderRepo.CreateSeckillOrder(ctx, order)
}

// 解析 Stream 消息：
func parseVoucherOrderMessage(msg redis.XMessage) (*model.VoucherOrder, error) {
	orderID, err := parseInt64FromStream(msg, "id")
	if err != nil {
		return nil, err
	}

	userID, err := parseInt64FromStream(msg, "userId")
	if err != nil {
		return nil, err
	}

	voucherID, err := parseInt64FromStream(msg, "voucherId")
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &model.VoucherOrder{
		ID:         orderID,
		UserID:     userID,
		VoucherID:  voucherID,
		Status:     1,
		CreateTime: now,
		UpdateTime: now,
	}, nil
}

// 工具函数：
func parseInt64FromStream(msg redis.XMessage, field string) (int64, error) {
	value, ok := msg.Values[field]
	if !ok {
		return 0, fmt.Errorf("stream field %s missing", field)
	}

	id, err := strconv.ParseInt(fmt.Sprint(value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("stream field %s invalid: %w", field, err)
	}

	return id, nil
}

// 处理 Pending List
func (s *voucherOrderService) handlePendingList(ctx context.Context) {
	//pending 消息失败超过 3 次后会进入 DLQ，而不是永远循环。
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		streams, err := s.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    constants.SeckillConsumerGroup,
			Consumer: constants.SeckillConsumerName,
			Streams:  []string{constants.SeckillStreamKey, "0"},
			Count:    1,
			Block:    0,
		}).Result()

		if errors.Is(err, redis.Nil) {
			return
		}

		if err != nil {
			log.Printf("[seckill] read pending failed: %v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if len(streams) == 0 || len(streams[0].Messages) == 0 {
			return
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				ok := s.handleMessageWithRetry(ctx, msg)
				if !ok {
					return
				}
			}
		}
	}
}

// 统一处理消息：成功 ACK，失败重试，超过次数进 DLQ
func (s *voucherOrderService) handleMessageWithRetry(ctx context.Context, msg redis.XMessage) bool {
	err := s.handleVoucherOrderMessage(ctx, msg)
	if err == nil {
		s.ackMessage(ctx, msg.ID)
		s.clearRetryCount(ctx, msg.ID)
		return true
	}

	retryCount, retryErr := s.increaseRetryCount(ctx, msg.ID)
	if retryErr != nil {
		log.Printf("[seckill] increase retry count failed, msgID=%s err=%v", msg.ID, retryErr)
		return false
	}

	if retryCount < constants.SeckillMaxRetry {
		log.Printf("[seckill] message will retry later, msgID=%s retry=%d", msg.ID, retryCount)
		return false
	}

	if dlqErr := s.moveToDLQ(ctx, msg, err, retryCount); dlqErr != nil {
		log.Printf("[seckill] move message to DLQ failed, msgID=%s err=%v", msg.ID, dlqErr)
		return false
	}

	s.ackMessage(ctx, msg.ID)
	s.clearRetryCount(ctx, msg.ID)
	return true
}

// ACK 封装
func (s *voucherOrderService) ackMessage(ctx context.Context, msgID string) {
	if err := s.redisClient.XAck(
		ctx,
		constants.SeckillStreamKey,
		constants.SeckillConsumerGroup,
		msgID,
	).Err(); err != nil {
		log.Printf("[seckill] ack message failed, msgID=%s err=%v", msgID, err)
	}
}

func (s *voucherOrderService) increaseRetryCount(ctx context.Context, msgID string) (int64, error) {
	key := constants.SeckillRetryKey + msgID

	count, err := s.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// 设置过期时间，避免 retry key 永久堆积。
	_ = s.redisClient.Expire(ctx, key, 24*time.Hour).Err()

	return count, nil
}

// 失败次数统计
func (s *voucherOrderService) clearRetryCount(ctx context.Context, msgID string) {
	key := constants.SeckillRetryKey + msgID
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		log.Printf("[seckill] clear retry count failed, msgID=%s err=%v", msgID, err)
	}
}

// 写入死信队列
func (s *voucherOrderService) moveToDLQ(ctx context.Context, msg redis.XMessage, cause error, retryCount int64) error {
	values := map[string]interface{}{
		"originMsgId": msg.ID,
		"reason":      cause.Error(),
		"retryCount":  retryCount,
		"failedAt":    time.Now().Format(time.RFC3339),
	}

	// 把原消息字段也带过去，方便后面人工排查或补偿。
	for k, v := range msg.Values {
		values[k] = fmt.Sprint(v)
	}

	return s.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: constants.SeckillDLQStreamKey,
		Values: values,
	}).Err()
}
