package tushare

import (
	"context"
	"encoding/json"
	"financia/public"
	"financia/public/db/model"
	"financia/util"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

func DailyStockAll(ctx context.Context, req *DailyReq) []*model.StockData {
	r := tuSharePost(public.TuShareDaily, req)
	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[DailyStockAll] [json.Marshal] [err] = %s", err.Error())
		return nil
	}

	var resp *DailyResp
	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[DailyStockAll] [json.Unmarshal] [err] = %s", err.Error())
		return nil
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

	return data
}

func DailyFundAll(ctx context.Context, req *DailyReq) []*model.FundData {
	r := tuSharePost(public.TuShareFundDaily, req)

	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[DailyFundAll] [json.Marshal] [err] = %s", err.Error())
		return nil
	}

	var resp *DailyResp
	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[DailyFundAll] [json.Unmarshal] [err] = %s", err.Error())
		return nil
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

	return data
}

func FundSalesRatio(ctx context.Context) []*FundSalesRatioResp {
	r := tuSharePost(public.TuShareFundSalesRatio, nil)

	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[FundSalesRatio] [json.Marshal] [err] = %s", err.Error())
		return nil
	}

	var resp *DailyResp
	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[FundSalesRatio] [json.Unmarshal] [err] = %s", err.Error())
		return nil
	}

	list := make([]*FundSalesRatioResp, 0, len(resp.Items))

	for _, item := range resp.Items {
		list = append(list, &FundSalesRatioResp{
			Year:      cast.ToString(item[0]),
			Bank:      cast.ToFloat64(item[1]),
			SecComp:   cast.ToFloat64(item[2]),
			FundComp:  cast.ToFloat64(item[3]),
			IndepComp: cast.ToFloat64(item[4]),
			Rests:     cast.ToFloat64(item[5]),
		})

	}

	return list
}

func FundSalesVol(ctx context.Context) []*FundSalesVolResp {
	r := tuSharePost(public.TuShareFundSalesVol, nil)

	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[FundSalesVol] [json.Marshal] [err] = %s", err.Error())
		return nil
	}

	var resp *DailyResp
	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[FundSalesVol] [json.Unmarshal] [err] = %s", err.Error())
		return nil
	}

	list := make([]*FundSalesVolResp, 0, len(resp.Items))
	for _, item := range resp.Items {
		list = append(list, &FundSalesVolResp{
			Year:      cast.ToString(item[0]),
			Quarter:   cast.ToString(item[1]),
			InstName:  cast.ToString(item[2]),
			FundScale: cast.ToFloat64(item[3]),
			Scale:     cast.ToFloat64(item[4]),
			Rank:      cast.ToInt(item[5]),
		})
	}

	return list
}

func FutTradeCal(ctx context.Context) ([]*FutTradeCalResp, []*FutTradeCalResp) {
	now := time.Now().Add(-31 * 24 * time.Hour).Format("20060102")
	end := time.Now().Add(31 * 24 * time.Hour).Format("20060102")
	r := tuSharePost(public.TuShareFutTradeCal, &DailyReq{
		Exchange:  "SSE",
		StartDate: now,
		EndDate:   end,
	})
	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[FundSalesVol] [json.Marshal] [err] = %s", err.Error())
		return nil, nil
	}

	var resp *DailyResp
	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[FundSalesVol] [json.Unmarshal] [err] = %s", err.Error())
		return nil, nil
	}

	var sse []*FutTradeCalResp
	for _, item := range resp.Items {
		sse = append(sse, &FutTradeCalResp{
			CalDate: cast.ToString(item[1]),
			IsOpen:  cast.ToInt(item[2]),
		})
	}

	r = tuSharePost(public.TuShareFutTradeCal, &DailyReq{
		Exchange:  "SZSE",
		StartDate: now,
		EndDate:   end,
	})
	marshal, err = json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[FundSalesVol] [json.Marshal] [err] = %s", err.Error())
		return nil, nil
	}

	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[FundSalesVol] [json.Unmarshal] [err] = %s", err.Error())
		return nil, nil
	}

	var szse []*FutTradeCalResp
	for _, item := range resp.Items {
		szse = append(szse, &FutTradeCalResp{
			CalDate: cast.ToString(item[1]),
			IsOpen:  cast.ToInt(item[2]),
		})
	}

	return sse, szse
}

func FutWeeklyDetail(ctx context.Context, prd string) []*FutWeeklyDetailResp {
	r := tuSharePost(public.TuShareFutWeeklyDetail, &DailyReq{
		Prd: prd,
	})

	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[FutWeeklyDetail] [json.Marshal] [err] = %s", err.Error())
		return nil
	}

	var resp *DailyResp
	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[FutWeeklyDetail] [json.Unmarshal] [err] = %s", err.Error())
		return nil
	}

	zap.S().Debugf("[FutWeeklyDetail] [resp] = %#v", resp.Items[0])
	list := make([]*FutWeeklyDetailResp, 0, len(resp.Items))
	for _, item := range resp.Items {
		weekDate := cast.ToString(item[16])
		if weekDate == "" {
			continue
		}
		list = append(list, &FutWeeklyDetailResp{
			Vol:          cast.ToInt(item[3]),
			VolYoy:       cast.ToFloat64(item[4]),
			Amount:       cast.ToFloat64(item[5]),
			AmountYoy:    cast.ToFloat64(item[6]),
			CumVol:       cast.ToInt(item[7]),
			CumVolYoy:    cast.ToFloat64(item[8]),
			Cumamt:       cast.ToFloat64(item[9]),
			CumamtYoy:    cast.ToFloat64(item[10]),
			OpenInterest: cast.ToInt(item[11]),
			InterestWow:  cast.ToFloat64(item[12]),
			McClose:      cast.ToFloat64(item[13]),
			CloseWow:     cast.ToFloat64(item[14]),
			WeekDate:     weekDate,
		})
	}

	return list
}
