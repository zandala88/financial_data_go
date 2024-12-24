package router

import (
	"errors"
	"financia/config"
	"financia/public/middlewares"
	"financia/public/vaildator"
	"financia/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

// HTTPRouter http 路由
func HTTPRouter() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

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
