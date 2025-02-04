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
		// 顶部tab
		auth.GET("/tab/list", common.GetTabList)

		// 公司 - 筛选参数
		auth.GET("/company/query", company.QueryCompany)
		// 公司 - 列表
		auth.GET("/company/list", company.ListCompany)
		// 公司 - 详情
		auth.GET("/company", company.DetailCompany)

		// 股票 - 筛选参数
		auth.GET("/stock/query", stock.QueryStock)
		// 股票 - 列表
		auth.GET("/stock/list", stock.ListStock)
		// 股票 - 判断是否有数据
		auth.GET("/stock/have", stock.HaveStock)
		// 股票 - 数据
		auth.GET("/stock/data", stock.DataStock)
		// 股票 - 详情中的信息
		auth.GET("/stock/info", stock.InfoStock)
		// 股票 - 首页图表
		auth.GET("/stock/graph", stock.GraphStock)
		// 股票 - 利润表
		auth.GET("/stock/income", stock.IncomeStock)
		// 股票 - 业绩预告
		auth.GET("/stock/forecast", stock.ForecastStock)
		// 股票 - 详情 - 十大股东
		auth.GET("/stock/top10", stock.Top10Stock)
		// 股票 - 首页排行
		auth.GET("/stock/hsgt/top10", stock.Top10HsgtStock)
		// 股票 - 预测数据
		auth.GET("/stock/predict", stock.PredictStock)

		// 公募基金 - 筛选参数
		auth.GET("/fund/query", fund.QueryFund)
		// 公募基金 - 列表
		auth.GET("/fund/list", fund.ListFund)
		// 公募基金 - 判断是否有数据
		auth.GET("/fund/have", fund.HaveFund)
		// 公募基金 - 数据
		auth.GET("/fund/data", fund.DataFund)
		// 公募基金 - 首页图表
		auth.GET("/fund/graph", fund.GraphFund)

		// 期货 - 筛选参数
		auth.GET("/fut/query", fut.QueryFut)
		// 期货 - 日历
		auth.GET("/fut/cal", fut.CalFut)
		// 期货 - 数据
		auth.GET("/fut/detail", fut.DetailFut)

		// 宏观经济 - shibor利率
		auth.GET("/economics/shibor", economics.ShiborEconomics)
		// 宏观经济 - GDP
		auth.GET("/economics/cn_gdp", economics.CnGdpEconomics)
		// 宏观经济 - CPI
		auth.GET("/economics/cn_cpi", economics.CnCpiEconomics)
	}

	httpAddr := fmt.Sprintf("%s:%s", config.Configs.App.IP, config.Configs.App.Port)
	if err := r.Run(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.S().Fatalf("listen: %s\n", err)
	}
}
