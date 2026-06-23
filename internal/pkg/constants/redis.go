package constants

// Redis key 前缀统一放在这里，避免项目里到处手写字符串。
//
// 例如登录 token 存 Redis 时会用:
// login:token:{token}
const (
	LoginCodeKey = "login:code:"
	LoginUserKey = "login:token:"

	CacheShopKey     = "cache:shop:"
	CacheShopTypeKey = "cache:shopType"
	LockShopKey      = "lock:shop:"

	SeckillStockKey = "seckill:stock:"
	BlogLikedKey    = "blog:liked:"
	FeedKey         = "feed:"
	ShopGeoKey      = "shop:geo:"
	UserSignKey     = "sign:"

	FollowsKey = "follows:"
)

// Redis 过期时间统一放在这里。
//
// 目前单位写在常量名里：Minutes 表示分钟，Seconds 表示秒。
const (
	LoginCodeTTLMinutes = 2
	LoginUserTTLMinutes = 30
	CacheNullTTLMinutes = 2
	CacheShopTTLMinutes = 30
	LockShopTTLSeconds  = 10
)
