package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/constants"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	errShopNotFound = errors.New("shop not found")
)

// shopCacheData 是写入 Redis 的店铺缓存包装对象。
//
// 这里不只是缓存店铺本身，还额外缓存一个 ExpireTime。
// 这个时间叫“逻辑过期时间”：Redis key 本身可以不过期，
// 但业务代码读出来后会判断 ExpireTime 是否已经过期。
type shopCacheData struct {
	Data       model.Shop `json:"data"`
	ExpireTime time.Time  `json:"expireTime"`
}

type ShopService interface {
	// QueryByID 查询店铺详情。
	QueryByID(ctx context.Context, id int64) result.Result
	// SaveShop 新增店铺。
	SaveShop(ctx context.Context, shop model.Shop) result.Result
	// UpdateShop 更新店铺。
	UpdateShop(ctx context.Context, shop model.Shop) result.Result
	// QueryByType 根据分类查询店铺列表。
	QueryByType(ctx context.Context, typeID int64, current int) result.Result
	// QueryByName 根据关键词搜索店铺。
	QueryByName(ctx context.Context, name string, current int) result.Result
}

type shopService struct {
	shopRepo    repository.ShopRepository // shopRepo 负责 tb_shop 数据库操作。
	redisClient *redis.Client             // redisClient 后面用于店铺缓存和 GEO 查询。
}

// NewShopService 创建店铺 Service。
func NewShopService(shopRepo repository.ShopRepository, redisClient *redis.Client) ShopService {
	return &shopService{shopRepo: shopRepo, redisClient: redisClient}
}

// QueryByID 查询店铺详情。
//
// 后面会重点学习 Redis 缓存穿透、缓存击穿、逻辑过期等内容。
func (s *shopService) QueryByID(ctx context.Context, id int64) result.Result {
	if id <= 0 {
		return result.Fail("invalid shop id")
	}

	shop, err := s.queryShopByIDWithCache(ctx, id)
	if errors.Is(err, errShopNotFound) {
		return result.Fail("shop not found")
	}
	if err != nil {
		log.Printf("[ShopService.QueryByID] query shop failed, id=%d, err=%v", id, err)
		return result.Fail("query shop failed")
	}
	return result.OKWithData(shop)
}

// SaveShop 新增店铺，保存成功后返回店铺 id。
func (s *shopService) SaveShop(ctx context.Context, shop model.Shop) result.Result {
	// TODO: Save shop to MySQL and return generated id.
	return result.Fail("TODO: save shop")
}

// UpdateShop 更新店铺。
//
// 更新 MySQL 后要删除 Redis 缓存，避免前端看到旧数据。
func (s *shopService) UpdateShop(ctx context.Context, shop model.Shop) result.Result {
	if shop.ID <= 0 {
		return result.Fail("invalid shop id")
	}

	if err := s.shopRepo.UpdateShop(ctx, &shop); err != nil {
		log.Printf("[ShopService.UpdateShop] update shop failed, id=%d, err=%v", shop.ID, err)
		return result.Fail("update shop failed")
	}

	// 先更新数据库，再删除缓存。
	// 这样下一次查询会重新从 MySQL 加载新数据，避免长期返回旧店铺信息。
	key := constants.CacheShopKey + strconv.FormatInt(shop.ID, 10)
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		log.Printf("[ShopService.UpdateShop] delete shop cache failed, key=%s, err=%v", key, err)
	}

	return result.OK()
}

// QueryByType 根据店铺分类分页查询。
func (s *shopService) QueryByType(ctx context.Context, typeID int64, current int) result.Result {
	if typeID < 0 {
		return result.Fail("invalid typeID")
	}
	shops, err := s.shopRepo.FindShopsByType(ctx, typeID, current)
	if err != nil {
		return result.Fail("quey shop by type failed")
	}
	return result.OKWithData(shops)
}

// QueryByName 根据店铺名称搜索。
func (s *shopService) QueryByName(ctx context.Context, name string, current int) result.Result {
	shops, err := s.shopRepo.FindShopsByName(ctx, name, current)
	if err != nil {
		return result.Fail("query shop by name failed")
	}
	return result.OKWithData(shops)
}

// queryShopByIDWithCache 是店铺详情缓存查询的核心入口。
//
// 这套逻辑同时处理三类问题：
// 1. 缓存穿透：不存在的店铺 id 会缓存空字符串，避免反复打 MySQL。
// 2. 缓存击穿：缓存未命中时用互斥锁，避免大量请求同时查 MySQL。
// 3. 热点重建：缓存逻辑过期后先返回旧数据，再后台异步重建缓存。
func (s *shopService) queryShopByIDWithCache(ctx context.Context, id int64) (*model.Shop, error) {
	shop, hit, expired, err := s.getShopFromCache(ctx, id)
	if err != nil {
		return nil, err
	}

	// Redis 有缓存，并且没有逻辑过期，直接返回。
	if hit && !expired {
		return shop, nil
	}

	// Redis 有缓存，但逻辑已经过期。
	// 此时仍然先返回旧数据，保证用户请求不被 MySQL 阻塞；
	// 同时只有抢到锁的请求负责后台重建缓存。
	if hit && expired {
		s.rebuildShopCacheAsync(id)
		return shop, nil
	}

	// Redis 完全没有缓存，说明可能是冷启动或 key 被删除。
	// 这种情况不能返回旧数据，只能走互斥锁同步重建。
	return s.queryShopOnCacheMiss(ctx, id)
}

// getShopFromCache 从 Redis 读取店铺缓存。
//
// 返回值含义：
// - shop：读取到的店铺数据。
// - hit：Redis 是否命中。即使命中的是空值，hit 也为 true。
// - expired：逻辑过期时间是否已过期。
// - err：读取或解析过程中出现的错误。
func (s *shopService) getShopFromCache(ctx context.Context, id int64) (*model.Shop, bool, bool, error) {
	key := constants.CacheShopKey + strconv.FormatInt(id, 10)

	cacheValue, err := s.redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, false, nil
	}
	if err != nil {
		return nil, false, false, err
	}

	// 命中空字符串，代表这个店铺在 MySQL 里不存在。
	// 这是“缓存空值”，专门用来解决缓存穿透。
	if cacheValue == "" {
		return nil, true, false, errShopNotFound
	}

	var cacheData shopCacheData
	if err := json.Unmarshal([]byte(cacheValue), &cacheData); err != nil {
		// 如果缓存格式坏了，删除这个 key，让下一次请求重新加载。
		log.Printf("[ShopService.getShopFromCache] invalid cache json, key=%s, err=%v", key, err)
		_ = s.redisClient.Del(ctx, key).Err()
		return nil, false, false, nil
	}

	expired := time.Now().After(cacheData.ExpireTime)
	return &cacheData.Data, true, expired, nil
}

// queryShopOnCacheMiss 处理 Redis 完全未命中的情况。
//
// 没命中时不能让所有请求都去查数据库，所以这里使用 Redis SETNX 做互斥锁。
// 抢到锁的请求负责查 MySQL 并写入缓存；没抢到锁的请求稍等一下再重试。
func (s *shopService) queryShopOnCacheMiss(ctx context.Context, id int64) (*model.Shop, error) {
	lockKey := constants.LockShopKey + strconv.FormatInt(id, 10)

	locked, err := s.tryLock(ctx, lockKey)
	if err != nil {
		return nil, err
	}

	if !locked {
		timer := time.NewTimer(50 * time.Millisecond)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
			return s.queryShopByIDWithCache(ctx, id)
		}
	}

	defer s.unlock(context.Background(), lockKey)

	// 双重检查：
	// 抢到锁后再查一次 Redis，避免“前一个请求刚写完缓存，当前请求又查一次 MySQL”。
	shop, hit, _, err := s.getShopFromCache(ctx, id)
	if err != nil {
		return nil, err
	}
	if hit {
		return shop, nil
	}

	return s.rebuildShopCache(ctx, id)
}

// rebuildShopCache 从 MySQL 查询店铺，并把结果重建到 Redis。
func (s *shopService) rebuildShopCache(ctx context.Context, id int64) (*model.Shop, error) {
	shop, err := s.shopRepo.FindShopByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据库也没有，缓存空字符串，避免同一个非法 id 反复打到 MySQL。
		if cacheErr := s.setShopNullCache(ctx, id); cacheErr != nil {
			log.Printf("[ShopService.rebuildShopCache] set null cache failed, id=%d, err=%v", id, cacheErr)
		}
		return nil, errShopNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := s.setShopWithLogicalExpire(ctx, *shop); err != nil {
		return nil, err
	}

	return shop, nil
}

// rebuildShopCacheAsync 在后台异步重建逻辑过期的店铺缓存。
//
// 注意：只有抢到锁的请求会开启 goroutine，避免并发重建。
func (s *shopService) rebuildShopCacheAsync(id int64) {
	lockKey := constants.LockShopKey + strconv.FormatInt(id, 10)

	locked, err := s.tryLock(context.Background(), lockKey)
	if err != nil || !locked {
		if err != nil {
			log.Printf("[ShopService.rebuildShopCacheAsync] try lock failed, id=%d, err=%v", id, err)
		}
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		defer s.unlock(context.Background(), lockKey)

		if _, err := s.rebuildShopCache(ctx, id); err != nil && !errors.Is(err, errShopNotFound) {
			log.Printf("[ShopService.rebuildShopCacheAsync] rebuild cache failed, id=%d, err=%v", id, err)
		}
	}()
}

// setShopWithLogicalExpire 把店铺写入 Redis，并附带逻辑过期时间。
func (s *shopService) setShopWithLogicalExpire(ctx context.Context, shop model.Shop) error {
	key := constants.CacheShopKey + strconv.FormatInt(shop.ID, 10)
	expireDuration := shopCacheLogicalExpireDuration(shop.ID)

	cacheData := shopCacheData{
		Data:       shop,
		ExpireTime: time.Now().Add(expireDuration),
	}

	cacheBytes, err := json.Marshal(cacheData)
	if err != nil {
		return err
	}

	// 逻辑过期模式下，Redis key 本身不设置 TTL。
	// 是否过期由 cacheData.ExpireTime 判断。
	return s.redisClient.Set(ctx, key, cacheBytes, 0).Err()
}

// setShopNullCache 缓存空值，专门防止缓存穿透。
func (s *shopService) setShopNullCache(ctx context.Context, id int64) error {
	key := constants.CacheShopKey + strconv.FormatInt(id, 10)
	return s.redisClient.Set(
		ctx,
		key,
		"",
		time.Duration(constants.CacheNullTTLMinutes)*time.Minute,
	).Err()
}

// tryLock 使用 Redis SETNX 尝试加锁。
func (s *shopService) tryLock(ctx context.Context, key string) (bool, error) {
	return s.redisClient.SetNX(
		ctx,
		key,
		"1",
		time.Duration(constants.LockShopTTLSeconds)*time.Second,
	).Result()
}

// unlock 释放 Redis 互斥锁。
func (s *shopService) unlock(ctx context.Context, key string) {
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		log.Printf("[ShopService.unlock] unlock failed, key=%s, err=%v", key, err)
	}
}

// shopCacheLogicalExpireDuration 返回店铺缓存的逻辑过期时长。
//
// id%10 是一个很轻量的过期时间扰动。
// 如果很多店铺同时写入缓存，这个扰动可以让它们的逻辑过期时间稍微错开，
// 避免大量热点 key 在同一分钟集中重建。
func shopCacheLogicalExpireDuration(id int64) time.Duration {
	base := time.Duration(constants.CacheShopTTLMinutes) * time.Minute
	jitter := time.Duration(id%10) * time.Minute
	return base + jitter
}
