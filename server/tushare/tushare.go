package tushare

import (
	"context"
	"financia/public"
	"financia/public/db/model"
	"financia/util"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

func DailyStockAll(ctx context.Context, req *DailyReq) []*model.StockData {
	r := tuSharePost(public.TuShareDaily, req, "")

	var resp DailyResp
	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[DailyStockAll] [marshalResp] [err] = %s", err.Error())
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
	r := tuSharePost(public.TuShareFundDaily, req, "")

	var resp DailyResp
	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[DailyFundAll] [marshalResp] [err] = %s", err.Error())
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
	r := tuSharePost(public.TuShareFundSalesRatio, nil, "")

	var resp DailyResp
	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[FundSalesRatio] [marshalResp] [err] = %s", err.Error())
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
	r := tuSharePost(public.TuShareFundSalesVol, nil, "")

	var resp DailyResp
	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[FundSalesVol] [marshalResp] [err] = %s", err.Error())
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
	}, "")

	var resp DailyResp
	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[FutTradeCal] [marshalResp] [err] = %s", err.Error())
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
	}, "")

	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[FutTradeCal] [marshalResp] [err] = %s", err.Error())
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
	}, "")

	var resp DailyResp
	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[FutWeeklyDetail] [marshalResp] [err] = %s", err.Error())
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

func StockIncome(ctx context.Context, tsCode string) []*StockIncomeResp {
	r := tuSharePost(public.TuShareStockIncome, &DailyReq{
		TsCode:     tsCode,
		ReportType: 1,
	}, "ann_date,basic_eps,total_revenue,total_cogs,"+
		"oper_exp,total_profit,income_tax,n_income,t_compr_income")

	var resp DailyResp
	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[StockIncome] [marshalResp] [err] = %s", err.Error())
		return nil
	}

	list := make([]*StockIncomeResp, 0, len(resp.Items))
	for _, item := range resp.Items {
		list = append(list, &StockIncomeResp{
			AnnDate:      util.ConvertDateStrToTime(cast.ToString(item[0]), timeLayout).Format(time.DateOnly),
			BasicEps:     cast.ToFloat64(item[1]),
			TotalRevenue: cast.ToFloat64(item[2]),
			TotalCogs:    cast.ToFloat64(item[3]),
			OperExp:      cast.ToFloat64(item[4]),
			TotalProfit:  cast.ToFloat64(item[5]),
			IncomeTax:    cast.ToFloat64(item[6]),
			NIncome:      cast.ToFloat64(item[7]),
			TComprIncome: cast.ToFloat64(item[8]),
		})
	}

	return list
}

func StockForecast(ctx context.Context, tsCode string) []*StockForecastResp {
	r := tuSharePost(public.TuShareStockForecast, &DailyReq{
		TsCode: tsCode,
	}, "ann_date,type,p_change_min,p_change_max,net_profit_min,"+
		"net_profit_max,last_parent_net,change_reason,update_flag")

	var resp DailyResp
	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[StockForecast] [marshalResp] [err] = %s", err.Error())
		return nil
	}

	list := make([]*StockForecastResp, 0, len(resp.Items))
	for _, item := range resp.Items {
		updateFlag := cast.ToString(item[8])
		if updateFlag == "1" {
			continue
		}

		list = append(list, &StockForecastResp{
			AnnDate:       util.ConvertDateStrToTime(cast.ToString(item[0]), timeLayout).Format(time.DateOnly),
			Type:          cast.ToString(item[1]),
			PChangeMin:    cast.ToFloat64(item[2]),
			PChangeMax:    cast.ToFloat64(item[3]),
			NetProfitMin:  cast.ToFloat64(item[4]),
			NetProfitMax:  cast.ToFloat64(item[5]),
			LastParentNet: cast.ToFloat64(item[6]),
			ChangeReason:  cast.ToString(item[7]),
		})
	}

	return list
}

func StockTop10(ctx context.Context, tsCode string) []*StockTop10Resp {
	r := tuSharePost(public.TuShareStockTop10, &DailyReq{
		TsCode: tsCode,
	}, "ann_date,holder_name,hold_amount,hold_ratio,hold_float_ratio,hold_change,holder_type")

	var resp DailyResp
	if err := marshalResp(r, &resp); err != nil {
		zap.S().Errorf("[StockTop10] [marshalResp] [err] = %s", err.Error())
		return nil
	}

	list := make([]*StockTop10Resp, 0, len(resp.Items))
	for _, item := range resp.Items {
		list = append(list, &StockTop10Resp{
			AnnDate:        util.ConvertDateStrToTime(cast.ToString(item[0]), timeLayout).Format(time.DateOnly),
			HolderName:     cast.ToString(item[1]),
			HoldAmount:     cast.ToFloat64(item[2]),
			HoldRatio:      cast.ToFloat64(item[3]),
			HoldFloatRatio: cast.ToFloat64(item[4]),
			HoldChange:     cast.ToFloat64(item[5]),
			HolderType:     cast.ToString(item[6]),
		})
	}

	return list
}
