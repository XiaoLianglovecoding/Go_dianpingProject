package idworker

// Worker 是分布式 ID 生成器的占位结构。
//
// Java 版黑马点评用 Redis 自增 + 时间戳生成订单 id；
// Go 版后面实现秒杀下单时，也会在这里补同样的能力。
type Worker struct {
	// TODO: Add Redis-backed sequence generation like the Java RedisIdWorker.
}

// NewWorker 创建一个 ID 生成器实例。
func NewWorker() *Worker {
	return &Worker{}
}

// NextID 根据业务前缀生成一个唯一 id。
//
// keyPrefix 可以是 order、voucher 等业务名。
func (w *Worker) NextID(keyPrefix string) (int64, error) {
	// TODO: Generate a globally unique ID with timestamp + Redis increment.
	return 0, nil
}
