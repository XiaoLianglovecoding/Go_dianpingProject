package constants

// 系统级常量，不属于某个具体业务表。
const (
	// UserNickNamePrefix 是新用户默认昵称前缀，例如 user_abcd1234。
	UserNickNamePrefix = "user_"
	DefaultPageSize    = 5  // 默认分页大小。
	MaxPageSize        = 10 // 前端列表常用分页大小。
)

const SeckillScript = `
-- 1. 获取参数
local voucherId = ARGV[1]
local userId = ARGV[2]
local orderId = ARGV[3]

-- 2. 数据 Key
local stockKey = 'seckill:stock:' .. voucherId   -- 存储库存的 String 结构
local orderKey = 'seckill:orders:' .. voucherId  -- 存储购买过该券的用户 ID 的 Set 结构

-- 3. 业务逻辑
-- 3.1 判断库存是否充足 (注意：如果 key 不存在默认当 0 处理)
if (tonumber(redis.call('get', stockKey)) or 0) <= 0 then
    return 1 -- 返回 1 代表库存不足
end

-- 3.2 判断是否已经下过单 (SISMEMBER 判断用户是否在 Set 中)
if (redis.call('sismember', orderKey, userId) == 1) then
    return 2 -- 返回 2 代表重复下单
end

-- 4. 扣库存，记录用户，发送消息到 MQ
-- 4.1 扣减库存
redis.call('decr', stockKey)
-- 4.2 将用户加入已下单集合
redis.call('sadd', orderKey, userId)
-- 4.3 发送消息到 Redis Stream 队列 (XADD stream.orders * userId xx voucherId xx id xx)
redis.call('xadd', 'stream.orders', '*', 'userId', userId, 'voucherId', voucherId, 'id', orderId)

return 0 -- 返回 0 代表抢购成功
`
