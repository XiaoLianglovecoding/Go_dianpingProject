package constants

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
)

const (
	LoginCodeTTLMinutes = 2
	LoginUserTTLMinutes = 30
	CacheNullTTLMinutes = 2
	CacheShopTTLMinutes = 30
	LockShopTTLSeconds  = 10
)
