package router

import (
	"go-base/handler"
	"go-base/middleware"
	"go-base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 初始化 service
	userService := service.NewUserService(db)

	// 初始化 handler
	userHandler := handler.NewUserHandler(userService)

	// 路由组 - 公开接口（不需要认证）
	public := r.Group("/api")
	{
		// 用户注册
		public.POST("/register", userHandler.Register)
		// 用户登录
		public.POST("/login", userHandler.Login)
	}

	// 路由组 - 需要认证的接口
	protected := r.Group("/api", middleware.JWTAuth())
	{
		// 获取当前登录用户信息
		protected.GET("/userInfo", userHandler.GetUserInfo)
	}

	return r
}
