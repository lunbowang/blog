package routes

import (
	"net/http"
	"time"
	"web_app/controller"

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	_ "web_app/docs" // 千万不要忘了导入把你上一步生成的docs
	"web_app/logger"
	"web_app/middlewares"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	//创建默认禁用中间件的Gin引擎
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	//使用自定义中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1 := r.Group("/api/v1")
	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	v1.GET("/posts", controller.GetPostListHandler)
	//根据时间或分数或获取帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2)
	v1.Use(middlewares.JWTAuthMiddleware()) //认证中间件

	{
		v1.POST("/post", controller.CreatePostHandler)
		//投票
		v1.POST("/vote", controller.PostVoteController)
	}
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	//r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	//	//如果是登录的用户，判断请求头中是否有效的JWT
	//	c.String(http.StatusOK, "pong")
	//})
	pprof.Register(r) //注册pprof相关路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
