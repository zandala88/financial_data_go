package router

import (
	"errors"
	"financia/config"
	"financia/public/middlewares"
	"financia/public/vaildator"
	"financia/service"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// HTTPRouter http 路由
func HTTPRouter() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                                     // 允许的来源，可以是单个或多个地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},                              // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                                            // 允许暴露的响应头
		AllowCredentials: true,                                                                  // 是否允许携带身份凭证（如 Cookie）
		MaxAge:           12 * time.Hour,                                                        // 浏览器预检请求的缓存时间
	})

	// 将 CORS 中间件应用于所有路由
	r.Use(corsMiddleware)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date", vaildator.DateValidator)
	}

	r.POST("/login", service.Login)
	r.POST("/register", service.Register)

	auth := r.Group("", middlewares.AuthCheck())
	{
		auth.GET("/info", service.GetInfo)

		auth.GET("/stock", service.GetStock)
	}

	httpAddr := fmt.Sprintf("%s:%s", config.Configs.App.IP, config.Configs.App.Port)
	if err := r.Run(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.S().Fatalf("listen: %s\n", err)
	}
}
