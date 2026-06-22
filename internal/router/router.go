package router

import (
	"net/http"
	"strings"

	"hmdp-go/internal/handler"
	"hmdp-go/internal/middleware"
	"hmdp-go/internal/pkg/result"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// NewRouter 创建 Gin 路由引擎，并注册所有接口。
//
// 可以把这里理解成 Spring MVC 里所有 @RequestMapping 的集中登记处。
func NewRouter(handlers *handler.Handlers, redisClient *redis.Client) *gin.Engine {
	// gin.Default() 自带日志和崩溃恢复中间件，适合开发阶段使用。
	r := gin.Default()
	// 全局挂载 token 刷新中间件；现在只是占位，不会拦截请求。
	r.Use(middleware.RefreshTokenMiddleware(redisClient))

	registerUserRoutes(r, handlers.User)
	registerShopRoutes(r, handlers.Shop)
	registerShopTypeRoutes(r, handlers.ShopType)
	registerBlogRoutes(r, handlers.Blog)
	registerFollowRoutes(r, handlers.Follow)
	registerUploadRoutes(r, handlers.Upload)
	registerVoucherRoutes(r, handlers.Voucher)
	registerVoucherOrderRoutes(r, handlers.VoucherOrder)

	return r
}

// registerUserRoutes 注册 /user 开头的用户接口。
func registerUserRoutes(r *gin.Engine, h *handler.UserHandler) {
	// ==========================================
	// 【公开 user 路由组】(不受 LoginInterceptor 拦截)
	// ==========================================
	publicGroup := r.Group("/user")
	{
		publicGroup.POST("/code", h.SendCode)
		publicGroup.POST("/login", h.Login)

		// 具体的静态路由优先
		publicGroup.GET("/info/:id", h.QueryUserInfo)
		// 泛化参数路由垫底
		publicGroup.GET("/:id", h.QueryUserByID)
	}

	// ==========================================
	// 【登录 user 路由组】(必须带有 Token 才能访问)
	// ==========================================
	protectedGroup := r.Group("/user")
	// 挂载登录拦截器 (保安)
	protectedGroup.Use(middleware.LoginInterceptor())
	{
		protectedGroup.GET("/me", h.Me)
		protectedGroup.POST("/logout", h.Logout)
		protectedGroup.POST("/sign", h.Sign)
		protectedGroup.GET("/sign/count", h.SignCount)
	}
}

// registerShopRoutes 注册 /shop 开头的店铺接口。
func registerShopRoutes(r *gin.Engine, h *handler.ShopHandler) {
	group := r.Group("/shop")
	group.POST("", h.SaveShop)
	group.PUT("", h.UpdateShop)
	group.GET("/of/type", h.QueryShopByType)
	group.GET("/of/name", h.QueryShopByName)
	group.GET("/:id", h.QueryShopByID)
}

// registerShopTypeRoutes 注册 /shop-type 开头的店铺分类接口。
func registerShopTypeRoutes(r *gin.Engine, h *handler.ShopTypeHandler) {
	group := r.Group("/shop-type")
	group.GET("/list", h.List)
}

// registerBlogRoutes 注册 /blog 开头的博客接口。
func registerBlogRoutes(r *gin.Engine, h *handler.BlogHandler) {
	// ==========================================
	// 【公开 blog 路由组】(不受 LoginInterceptor 拦截)
	// ==========================================
	publicgroup := r.Group("/blog")
	{
		publicgroup.GET("/hot", h.QueryHotBlog)
		publicgroup.GET("/likes/:id", h.QueryBlogLikes)
		publicgroup.GET("/of/user", h.QueryBlogByUserID)
		publicgroup.GET("/:id", h.QueryByID)
	}

	// ==========================================
	// 【登录 blog 路由组】(必须带有 Token 才能访问)
	// ==========================================
	protectedGroup := r.Group("/blog")
	// 挂载登录拦截器 (保安)
	protectedGroup.Use(middleware.LoginInterceptor())
	{
		protectedGroup.GET("/of/me", h.QueryMyBlog)
		protectedGroup.GET("/of/follow", todoRoute("TODO: query blog of follow"))
		protectedGroup.POST("", h.SaveBlog)
		protectedGroup.PUT("/like/:id", h.LikeBlog)
	}
}

// registerFollowRoutes 注册 /follow 开头的关注接口。
func registerFollowRoutes(r *gin.Engine, h *handler.FollowHandler) {
	group := r.Group("/follow")
	group.Use(middleware.LoginInterceptor())
	{
		group.GET("/or/not/:id", h.IsFollow)
		group.PUT("/:id/:isFollow", h.Follow)
		group.GET("/common/:id", h.Common)
	}
}

// registerUploadRoutes 注册 /upload 开头的上传接口。
func registerUploadRoutes(r *gin.Engine, h *handler.UploadHandler) {
	group := r.Group("/upload")
	group.POST("/blog", h.UploadBlog)
	group.GET("/blog/delete", h.DeleteBlog)
}

// registerVoucherRoutes 注册 /voucher 开头的优惠券接口。
func registerVoucherRoutes(r *gin.Engine, h *handler.VoucherHandler) {
	group := r.Group("/voucher")
	group.POST("", h.AddVoucher)
	group.POST("/seckill", h.AddSeckillVoucher)
	group.GET("/list/:shopId", h.QueryVoucherOfShop)
}

// registerVoucherOrderRoutes 注册 /voucher-order 开头的订单接口。
func registerVoucherOrderRoutes(r *gin.Engine, h *handler.VoucherOrderHandler) {
	group := r.Group("/voucher-order")
	group.POST("/seckill/:id", h.SeckillVoucher)
}

// todoRoute 用于临时占位尚未实现的接口。
//
// 这样前端访问时不会 404，而是能明确看到哪个接口还没写。
func todoRoute(message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, result.Fail(strings.TrimSpace(message)))
	}
}
