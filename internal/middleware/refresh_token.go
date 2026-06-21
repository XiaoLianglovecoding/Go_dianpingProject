package middleware

import (
	"hmdp-go/internal/dto"
	"hmdp-go/internal/pkg/constants"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RefreshTokenMiddleware 是 token 刷新中间件。
//
// Java 版有类似 RefreshTokenInterceptor：
// 如果请求带了合法 token，就刷新 Redis 里的 token 过期时间，
// 并把用户信息放到请求上下文，后续 Handler/Service 可以取。
func RefreshTokenMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 token (前端通常把 token 放在 headers 的 authorization 字段)
		token := c.GetHeader("authorization")
		if token == "" {
			// 如果没有 token，说明是未登录的游客，直接放行（后面如果有 LoginInterceptor 会拦截他）
			c.Next()
			return
		}

		// 2. 拼接 Redis Key
		key := constants.LoginUserKey + token

		// 3. 从 Redis 获取完整的 Hash 数据 (Java 里的 entries)
		// 注意：go-redis 的 HGetAll 如果查不到 key，会返回一个空 map，不会报错 redis.Nil，这极其方便！
		userMap, err := redisClient.HGetAll(c, key).Result()
		if err != nil || len(userMap) == 0 {
			// 查不到数据，说明 token 是伪造的，或者已经过期被清理了。直接放行。
			c.Next()
			return
		}

		// 4. 将查到的 Map 转换成对象 (UserDTO)
		// 因为我们在存的时候把 id 转成了 string，现在取出来要转回 int64
		// 【优化点 1】严谨处理 id 解析错误。一旦失败，说明 Redis 数据异常，直接当做无效 token 丢弃。
		id, err := strconv.ParseInt(userMap["id"], 10, 64)
		if err != nil {
			log.Printf("⚠️ 警告: Token (%s) 中的 ID 解析失败，Redis 脏数据? err: %v", key, err)
			c.Next()
			return
		}
		userDTO := dto.UserDTO{
			ID:       id,
			NickName: userMap["nickName"],
			Icon:     userMap["icon"],
		}

		// 5. 【核心魔法】把用户信息存入 Gin 的 Context 中！
		// 它的作用就跟 Java 里的 ThreadLocal 一模一样，供后续的 Handler 和 Service 读取。
		c.Set("user", userDTO)

		// 6. 刷新 Token 的有效期 (比如 30 分钟)
		// 这样只要用户一直在操作 App，他的登录状态就不会掉线。
		ttl := time.Duration(constants.LoginUserTTLMinutes) * time.Minute
		// 【优化点 2】捕捉 Expire 的错误并打印日志，不阻断请求。
		if err := redisClient.Expire(c, key, ttl).Err(); err != nil {
			log.Printf("⚠️ 警告: 刷新 Token (%s) 续期失败: %v", key, err)
		}

		// 7. 所有的校验和续期都搞定了，继续执行下一个拦截器或最终的业务逻辑
		c.Next()
	}
}
