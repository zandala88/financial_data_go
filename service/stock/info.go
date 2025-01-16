package stock

import "github.com/gin-gonic/gin"

type InfoStockReq struct {
	Id int `form:"id" binding:"required"`
}

type InfoStockResp struct {
}

func InfoStock(c *gin.Context) {
	// todo

}
