package common

import (
	"financia/util"
	"github.com/gin-gonic/gin"
)

type TabListResp struct {
	List []string `json:"list"`
}

func GetTabList(c *gin.Context) {
	resp := TabListResp{
		List: []string{"上市公司", "股票市场", "公募基金", "期货数据", "新闻快讯"},
	}
	util.SuccessResp(c, resp)
}
