package user

import (
	"context"
	"errors"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/dao"
	"financia/public/db/model"
	"financia/server/python"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
)

func predict(c context.Context, userId int64, resp *UserInfoResp) error {
	rdb := connector.GetRedis().WithContext(c)
	eg, ctx := errgroup.WithContext(c)
	stockIdList, fundIdList, err := dao.GetFollowList(ctx, userId)
	if err != nil {
		zap.S().Error("[Info] [GetFollowList] [err] = ", err.Error())
		return nil
	}

	eg.Go(func() error {
		if len(stockIdList) == 0 {
			return nil
		}

		stockList, stocksToPredict, err := stockPredictKeys(rdb, stockIdList, ctx)
		if err != nil {
			zap.S().Error("[Info] [stockPredictKeys] [err] = ", err.Error())
			return err
		}
		resp.StockList = stockList

		// 并发调用 Python 预测
		var mu sync.Mutex
		eg2, _ := errgroup.WithContext(ctx)

		for _, stock := range stocksToPredict {
			stock := stock // 避免闭包问题
			eg2.Go(func() error {
				// 获取最近 30 天的股票数据
				stockData, err := dao.GetStockDataLimit30(ctx, stock.TsCode)
				if err != nil {
					zap.S().Error("[Info] [GetStockDataLimit30] [err] = ", err.Error())
					return err
				}

				if len(stockData) != 31 {
					return nil
				}

				// 排序数据，确保时间顺序正确
				sort.Slice(stockData, func(i, j int) bool {
					return stockData[i].TradeDate.Before(stockData[j].TradeDate)
				})

				// 调用 Python 进行预测
				val, err := python.PythonPredictStock(stock.Id, stockData)
				if err != nil {
					zap.S().Error("[PredictStock] [PythonPredictStock] [err] = ", err.Error())
					return err
				}

				// 加锁，确保多线程安全
				mu.Lock()
				resp.StockList = append(resp.StockList, &UserInfoData{
					Id:      stock.Id,
					Name:    stock.Name,
					Val:     stockData[len(stockData)-1].Close,
					NextVal: val,
				})
				mu.Unlock()

				return nil
			})
		}

		// 等待所有 Python 预测任务完成
		if err := eg2.Wait(); err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		if len(fundIdList) == 0 {
			return nil
		}

		fundList, fundsToPredict, err := fundPredictKeys(rdb, fundIdList, ctx)
		if err != nil {
			zap.S().Error("[Info] [fundPredictKeys] [err] = ", err.Error())
			return err
		}
		resp.FundList = fundList

		var mu sync.Mutex
		eg2, _ := errgroup.WithContext(ctx)

		for _, fund := range fundsToPredict {
			fund := fund
			eg2.Go(func() error {
				fundData, err := dao.GetFundDataLimit30(ctx, fund.TsCode)
				if err != nil {
					zap.S().Error("[Info] [GetFundDataLimit30] [err] = ", err.Error())
					return err
				}

				if len(fundData) != 31 {
					return nil
				}

				sort.Slice(fundData, func(i, j int) bool {
					return fundData[i].TradeDate.Before(fundData[j].TradeDate)
				})

				val, err := python.PythonPredictFund(int(fund.Id), fundData)
				if err != nil {
					zap.S().Error("[PredictFund] [PythonPredictFund] [err] = ", err.Error())
					return err
				}

				mu.Lock()
				resp.FundList = append(resp.FundList, &UserInfoData{
					Id:      int(fund.Id),
					Name:    fund.Name,
					Val:     fundData[len(fundData)-1].Close,
					NextVal: val,
				})
				mu.Unlock()

				return nil
			})
		}

		if err := eg2.Wait(); err != nil {
			return err
		}

		return nil
	})

	// 等待所有任务完成
	if err := eg.Wait(); err != nil {
		zap.S().Error("[Info] [err] = ", err.Error())
		return err
	}

	return nil
}

func stockPredictKeys(rdb *redis.Client, stockIdList []int, ctx context.Context) ([]*UserInfoData, []*model.StockInfo, error) {
	stockInfos, err := dao.GetStockInfos(ctx, stockIdList)
	if err != nil {
		zap.S().Error("[Info] [GetStockInfos] [err] = ", err.Error())
		return nil, nil, err
	}

	// 组合 Redis Key
	predictKeys := make([]string, 0, len(stockInfos))
	for _, v := range stockInfos {
		predictKeys = append(predictKeys, fmt.Sprintf(public.RedisKeyStockPredict, v.Id))
	}
	pipe := rdb.Pipeline()
	predictCmd := pipe.MGet(ctx, predictKeys...)
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		zap.S().Error("[Info] [Pipeline] [err] = ", err.Error())
		return nil, nil, err
	}
	predictResults, _ := predictCmd.Result()

	var (
		stocksToPredict []*model.StockInfo
		stockList       []*UserInfoData
	)

	for i, v := range stockInfos {
		if predictResults[i] != nil {
			result, _ := rdb.Get(ctx, fmt.Sprintf(public.RedisKeyStockToday, v.TsCode)).Result()
			// 直接使用 Redis 预测值
			stockList = append(stockList, &UserInfoData{
				Id:      v.Id,
				Name:    v.Name,
				Val:     cast.ToFloat64(result),
				NextVal: cast.ToFloat64(predictResults[i]), // Redis 缓存预测值
			})
		} else {
			// 需要进行 Python 预测
			stocksToPredict = append(stocksToPredict, v)
		}
	}

	return stockList, stocksToPredict, nil
}

func fundPredictKeys(rdb *redis.Client, fundIdList []int, ctx context.Context) ([]*UserInfoData, []*model.FundInfo, error) {
	fundInfos, err := dao.GetFundInfos(ctx, fundIdList)
	if err != nil {
		zap.S().Error("[Info] [GetFundInfos] [err] = ", err.Error())
		return nil, nil, err
	}

	predictKeys := make([]string, 0, len(fundInfos))
	for _, v := range fundInfos {
		predictKeys = append(predictKeys, fmt.Sprintf(public.RedisKeyFundPredict, v.Id))
	}
	pipe := rdb.Pipeline()
	predictCmd := pipe.MGet(ctx, predictKeys...)
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		zap.S().Error("[Info] [Pipeline] [err] = ", err.Error())
		return nil, nil, err
	}
	predictResults, _ := predictCmd.Result()

	var (
		fundsToPredict []*model.FundInfo
		fundList       []*UserInfoData
	)

	for i, v := range fundInfos {
		if predictResults[i] != nil {
			result, _ := rdb.Get(ctx, fmt.Sprintf(public.RedisKeyFundToday, v.TsCode)).Result()
			fundList = append(fundList, &UserInfoData{
				Id:      int(v.Id),
				Name:    v.Name,
				Val:     cast.ToFloat64(result),
				NextVal: cast.ToFloat64(predictResults[i]),
			})
		} else {
			fundsToPredict = append(fundsToPredict, v)
		}
	}
	return fundList, fundsToPredict, nil
}
