package tushare

import (
	"context"
	"encoding/json"
	"financia/public/db/dao"
	"financia/public/db/model"
	"financia/util"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type DailyReq struct {
	TsCode    string `json:"ts_code,omitempty"`
	TradeDate string `json:"trade_date,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
}

type DailyResp struct {
	Fields []string        `json:"fields"`
	Items  [][]interface{} `json:"items"`
}

func DailyStockAll(ctx context.Context, tsCode string) bool {
	r := tuSharePost("daily", &DailyReq{
		TsCode: tsCode,
	})
	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[Daily] [json.Marshal] [err] = %s", err.Error())
		return false
	}

	var resp *DailyResp
	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[Daily] [json.Unmarshal] [err] = %s", err.Error())
		return false
	}

	data := make([]*model.StockData, 0, len(resp.Items))
	for _, item := range resp.Items {
		data = append(data, &model.StockData{
			TsCode:    cast.ToString(item[0]),
			TradeDate: util.ConvertDateStrToTime(cast.ToString(item[1]), timeLayout),
			Open:      cast.ToFloat64(item[2]),
			High:      cast.ToFloat64(item[3]),
			Low:       cast.ToFloat64(item[4]),
			Close:     cast.ToFloat64(item[5]),
			PreClose:  cast.ToFloat64(item[6]),
			Change:    cast.ToFloat64(item[7]),
			PctChg:    cast.ToFloat64(item[8]),
			Vol:       cast.ToInt64(item[9]),
			Amount:    cast.ToFloat64(item[10]),
		})
	}
	if err := dao.InsertStockData(ctx, data); err != nil {
		zap.S().Errorf("[Daily] [InsertStockData] [err] = %s", err.Error())
		return false
	}
	return len(data) > 0
}

func DailyFundAll(ctx context.Context, tsCode string) bool {
	r := tuSharePost("fund_daily", &DailyReq{
		TsCode: tsCode,
	})

	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[Daily] [json.Marshal] [err] = %s", err.Error())
		return false
	}

	var resp *DailyResp
	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[Daily] [json.Unmarshal] [err] = %s", err.Error())
		return false
	}

	data := make([]*model.FundData, 0, len(resp.Items))
	for _, item := range resp.Items {
		data = append(data, &model.FundData{
			TsCode:    cast.ToString(item[0]),
			TradeDate: util.ConvertDateStrToTime(cast.ToString(item[1]), timeLayout),
			Open:      cast.ToFloat64(item[2]),
			High:      cast.ToFloat64(item[3]),
			Low:       cast.ToFloat64(item[4]),
			Close:     cast.ToFloat64(item[5]),
			PreClose:  cast.ToFloat64(item[6]),
			Change:    cast.ToFloat64(item[7]),
			PctChg:    cast.ToFloat64(item[8]),
			Vol:       cast.ToInt64(item[9]),
			Amount:    cast.ToFloat64(item[10]),
		})
	}
	if err := dao.InsertFundData(ctx, data); err != nil {
		zap.S().Errorf("[Daily] [InsertFundData] [err] = %s", err.Error())
		return false
	}
	return len(data) > 0
}

// 每日执行获取当日数据
func Daily(ctx context.Context) {
	zap.S().Debugf("[Daily] [start]")
	r := tuSharePost("daily", &DailyReq{})
	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[Daily] [json.Marshal] [err] = %s", err.Error())
		return
	}

	var resp *DailyResp
	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[Daily] [json.Unmarshal] [err] = %s", err.Error())
		return
	}

	data := make([]*model.StockData, 0, len(resp.Items))

	// 获取库中所有stock，f_ts_code 转map
	stockList, _, err := dao.GetStockList(ctx, "", nil, nil, nil, 1, 10000)
	hash := map[string]struct{}{}
	for _, stock := range stockList {
		hash[stock.TsCode] = struct{}{}
	}

	for _, item := range resp.Items {
		tsCode := cast.ToString(item[0])
		if _, ok := hash[tsCode]; !ok {
			continue
		}

		data = append(data, &model.StockData{
			TsCode:    tsCode,
			TradeDate: util.ConvertDateStrToTime(cast.ToString(item[1]), timeLayout),
			Open:      cast.ToFloat64(item[2]),
			High:      cast.ToFloat64(item[3]),
			Low:       cast.ToFloat64(item[4]),
			Close:     cast.ToFloat64(item[5]),
			PreClose:  cast.ToFloat64(item[6]),
			Change:    cast.ToFloat64(item[7]),
			PctChg:    cast.ToFloat64(item[8]),
			Vol:       cast.ToInt64(item[9]),
			Amount:    cast.ToFloat64(item[10]),
		})
	}

	if err := dao.InsertStockData(ctx, data); err != nil {
		zap.S().Errorf("[Daily] [InsertStockData] [err] = %s", err.Error())
		return
	}
}
