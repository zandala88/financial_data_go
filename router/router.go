package router

import (
	"errors"
	"financia/config"
	"financia/public/vaildator"
	"financia/service/common"
	"financia/service/company"
	"financia/service/economics"
	"financia/service/fund"
	"financia/service/fut"
	"financia/service/stock"

	"financia/service/user"
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
		AllowOrigins:     []string{"*"},                                       // 允许的来源，可以是单个或多个地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"*"},                                       // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 允许暴露的响应头
		AllowCredentials: true,                                                // 是否允许携带身份凭证（如 Cookie）
		MaxAge:           12 * time.Hour,                                      // 浏览器预检请求的缓存时间
	})

	// 将 CORS 中间件应用于所有路由
	r.Use(corsMiddleware)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date", vaildator.DateValidator)
		v.RegisterValidation("email", vaildator.EmailValidator)
	}

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", user.Login)
		v1.POST("/register", user.Register)
		v1.GET("/code", user.Code)
	}

	//auth := v1.Use(middleware.AuthCheck())
	auth := v1.Use()
	{
		auth.GET("/tab/list", common.GetTabList)

		auth.GET("/company/query", company.QueryCompany)
		auth.GET("/company/list", company.ListCompany)
		auth.GET("/company", company.DetailCompany)

		auth.GET("/stock/query", stock.QueryStock)
		auth.GET("/stock/list", stock.ListStock)
		auth.GET("/stock/have", stock.HaveStock)
		auth.GET("/stock/data", stock.DataStock)
		auth.GET("/stock/info", stock.InfoStock)
		auth.GET("/stock/graph", stock.GraphStock)
		auth.GET("/stock/income", stock.IncomeStock)
		auth.GET("/stock/forecast", stock.ForecastStock)
		auth.GET("/stock/top10", stock.Top10Stock)
		auth.GET("/stock/hsgt/top10", stock.Top10HsgtStock)

		auth.GET("/fund/query", fund.QueryFund)
		auth.GET("/fund/list", fund.ListFund)
		auth.GET("/fund/have", fund.HaveFund)
		auth.GET("/fund/data", fund.DataFund)
		auth.GET("/fund/graph", fund.GraphFund)

		auth.GET("/fut/query", fut.QueryFut)
		auth.GET("/fut/cal", fut.CalFut)
		auth.GET("/fut/detail", fut.DetailFut)

		auth.GET("/economics/shibor", economics.ShiborEconomics)
		auth.GET("/economics/cn_gdp", economics.CnGdpEconomics)
		auth.GET("/economics/cn_cpi", economics.CnCpiEconomics)
	}

	httpAddr := fmt.Sprintf("%s:%s", config.Configs.App.IP, config.Configs.App.Port)
	if err := r.Run(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.S().Fatalf("listen: %s\n", err)
	}
}
