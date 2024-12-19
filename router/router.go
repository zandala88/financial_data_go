package router

import (
	"errors"
	"financia/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// HTTPRouter http 路由
func HTTPRouter() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	httpAddr := fmt.Sprintf("%s:%s", config.Configs.App.IP, config.Configs.App.HTTPServerPort)
	if err := r.Run(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.S().Fatalf("listen: %s\n", err)
	}
}
