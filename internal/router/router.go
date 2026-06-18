package router

import (
	"net/http"
	"strings"

	"hmdp-go/internal/handler"
	"hmdp-go/internal/middleware"
	"hmdp-go/internal/pkg/result"

	"github.com/gin-gonic/gin"
)

func NewRouter(handlers *handler.Handlers) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.RefreshTokenMiddleware())

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

func registerUserRoutes(r *gin.Engine, h *handler.UserHandler) {
	group := r.Group("/user")
	group.POST("/code", h.SendCode)
	group.POST("/login", h.Login)
	group.POST("/logout", h.Logout)
	group.POST("/sign", h.Sign)

	group.GET("/me", h.Me)
	group.GET("/sign/count", h.SignCount)
	group.GET("/info/:id", h.QueryUserInfo)
	group.GET("/:id", h.QueryUserByID)
}

func registerShopRoutes(r *gin.Engine, h *handler.ShopHandler) {
	group := r.Group("/shop")
	group.POST("", h.SaveShop)
	group.PUT("", h.UpdateShop)
	group.GET("/of/type", h.QueryShopByType)
	group.GET("/of/name", h.QueryShopByName)
	group.GET("/:id", h.QueryShopByID)
}

func registerShopTypeRoutes(r *gin.Engine, h *handler.ShopTypeHandler) {
	group := r.Group("/shop-type")
	group.GET("/list", h.List)
}

func registerBlogRoutes(r *gin.Engine, h *handler.BlogHandler) {
	group := r.Group("/blog")
	group.POST("", h.SaveBlog)
	group.PUT("/like/:id", h.LikeBlog)
	group.GET("/of/me", h.QueryMyBlog)
	group.GET("/hot", h.QueryHotBlog)
	group.GET("/likes/:id", h.QueryBlogLikes)
	group.GET("/of/user", h.QueryBlogByUserID)
	group.GET("/of/follow", todoRoute("TODO: query blog of follow"))
	group.GET("/:id", h.QueryByID)
}

func registerFollowRoutes(r *gin.Engine, h *handler.FollowHandler) {
	group := r.Group("/follow")
	group.GET("/or/not/:id", h.IsFollow)
	group.PUT("/:id/:isFollow", h.Follow)
	group.GET("/common/:id", h.Common)
}

func registerUploadRoutes(r *gin.Engine, h *handler.UploadHandler) {
	group := r.Group("/upload")
	group.POST("/blog", h.UploadBlog)
	group.GET("/blog/delete", h.DeleteBlog)
}

func registerVoucherRoutes(r *gin.Engine, h *handler.VoucherHandler) {
	group := r.Group("/voucher")
	group.POST("", h.AddVoucher)
	group.POST("/seckill", h.AddSeckillVoucher)
	group.GET("/list/:shopId", h.QueryVoucherOfShop)
}

func registerVoucherOrderRoutes(r *gin.Engine, h *handler.VoucherOrderHandler) {
	group := r.Group("/voucher-order")
	group.POST("/seckill/:id", h.SeckillVoucher)
}

func todoRoute(message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, result.Fail(strings.TrimSpace(message)))
	}
}
