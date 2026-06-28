package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	// 定义一个起始时间戳 (例如 2024-01-01 00:00:00 的毫秒数)
	// 这样可以把 41 位时间戳的寿命往后推 69 年
	epoch int64 = 1704067200000

	nodeBits     uint8 = 10 // 机器 ID 占用的位数
	sequenceBits uint8 = 12 // 序列号占用的位数

	// 计算最大值，用于校验 (-1 向左移位然后再异或)
	nodeMax      int64 = -1 ^ (-1 << nodeBits)     // 1023
	sequenceMask int64 = -1 ^ (-1 << sequenceBits) // 4095

	// 定义向左位移的偏移量
	nodeShift uint8 = sequenceBits            // 机器 ID 需要向左移 12 位
	timeShift uint8 = sequenceBits + nodeBits // 时间戳需要向左移 22 位 (12+10)
)

// Snowflake 雪花算法生成器结构体
type Snowflake struct {
	mu        sync.Mutex // Go 原生互斥锁，保证并发安全
	timestamp int64      // 上次生成 ID 的时间戳
	nodeID    int64      // 当前机器/节点 ID
	sequence  int64      // 当前毫秒内的序列号
}

// 全局唯一的雪花算法实例
var GlobalSnowflake *Snowflake

// InitSnowflake 初始化全局实例（在程序启动时调用）
func InitSnowflake(nodeID int64) error {
	if nodeID < 0 || nodeID > nodeMax {
		return errors.New("机器 NodeID 超出范围 (0-1023)")
	}
	GlobalSnowflake = &Snowflake{
		nodeID: nodeID,
	}
	return nil
}

// NextVal 生成下一个唯一 ID
func (s *Snowflake) NextVal() int64 {
	// 加锁，保证同一时刻只能有一个 Goroutine 进来生成 ID
	s.mu.Lock()
	defer s.mu.Unlock() // 离开函数时自动解锁

	now := time.Now().UnixMilli()

	if now == s.timestamp {
		// 如果在同一毫秒内，序列号自增 (+1)，并用掩码防止溢出
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果同一毫秒内的 4096 个序列号用完了，就死循环等待下一毫秒
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 如果进入了新的一毫秒，序列号重置为 0
		s.sequence = 0
	}

	// 记录本次生成的时间戳
	s.timestamp = now

	// 将时间戳、机器 ID、序列号通过位运算 (左移 << 和按位或 |) 拼装成 64 位整数
	id := ((now - epoch) << timeShift) |
		(s.nodeID << nodeShift) |
		(s.sequence)

	return id
}
