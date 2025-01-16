package fund

import (
	"financia/public/db/dao"
	"financia/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ListFundReq struct {
	Search     string   `form:"search"`
	FundType   []string `form:"fundType"`
	InvestType []string `form:"investType"`
	Page       int      `form:"page" binding:"required"`
	PageSize   int      `form:"pageSize" binding:"required"`
}

type ListFundResp struct {
	List         []*ListFundSimple `json:"list"`
	TotalPageNum int               `json:"totalPageNum"`
	HasMore      bool              `json:"hasMore"`
}

type ListFundSimple struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Management   string  `json:"management"`
	Custodian    string  `json:"custodian"`
	FundType     string  `json:"fundType"`
	IssueAmount  float64 `json:"issueAmount"`
	MFree        float64 `json:"mFree"`
	CFree        float64 `json:"cFree"`
	DurationYear float64 `json:"durationYear"`
	PValue       float64 `json:"pValue"`
	MinAmount    float64 `json:"minAmount"`
	ExpReturn    float64 `json:"expReturn"`
	Benchmark    string  `json:"benchmark"`
	InvestType   string  `json:"investType"`
	Type         string  `json:"type"`
	Trustee      string  `json:"trustee"`
}

func ListFund(c *gin.Context) {
	var req ListFundReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCode(c, util.ShouldBindJSONError)
		zap.S().Error("[ListStock] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	list, count, err := dao.GetFundList(c, req.Search, req.FundType, req.InvestType, req.Page, req.PageSize)
	if err != nil {
		util.FailRespWithCode(c, util.InternalServerError)
		zap.S().Error("[ListStock] [GetStockList] [err] = ", err.Error())
		return
	}

	respList := make([]*ListFundSimple, 0, len(list))
	for _, v := range list {
		respList = append(respList, &ListFundSimple{
			Id:           v.Id,
			Name:         v.Name,
			Management:   v.Management,
			Custodian:    v.Custodian,
			FundType:     v.FundType,
			IssueAmount:  v.IssueAmount,
			MFree:        v.MFee,
			CFree:        v.CCFee,
			DurationYear: v.DurationYear,
			PValue:       v.PValue,
			MinAmount:    v.MinAmount,
			ExpReturn:    v.ExpReturn,
			Benchmark:    v.Benchmark,
			InvestType:   v.InvestType,
			Type:         v.Type,
			Trustee:      v.Trustee,
		})
	}

	util.SuccessResp(c, &ListFundResp{
		List:         respList,
		HasMore:      count > int64(req.Page*(req.PageSize-1)+len(list)),
		TotalPageNum: int(count/int64(req.PageSize) + 1),
	})
}
