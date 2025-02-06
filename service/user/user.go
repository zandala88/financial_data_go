package user

import (
	"errors"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/dao"
	"financia/public/db/model"
	"financia/server"
	"financia/server/python"
	"financia/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
)

// Code 获取验证码
func Code(c *gin.Context) {
	var req GetCodeReq
	if err := c.ShouldBindQuery(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[Code] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	code, err := dao.GetEmailCode(c, req.Email)
	if err == nil || code != "" {
		util.FailRespWithCodeAndZap(c, util.CodeLimitError, "[Code] [GetEmailCode] [err] = ", "验证码发送过于频繁")
		return
	}

	code = public.GenerateVerificationCode(6)

	go server.SendEmail(req.Email, public.EmailTitle, code)

	err = dao.SetEmailCode(c, req.Email, code)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Code] [SetEmailCode] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, nil)
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[Login] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	userId := dao.GetUserId(c, req.Email)
	if userId <= 0 {
		util.FailRespWithCodeAndZap(c, util.ReqDataError, "[Login] [GetUserId] [err] = ", "用户不存在")
		return
	}

	token, err := util.GenerateJWT(userId)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Login] [GenerateJWT] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, LoginResp{Token: token})
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		util.FailRespWithCodeAndZap(c, util.ShouldBindJSONError, "[Register] [ShouldBindJSON] [err] = ", err.Error())
		return
	}

	code, err := dao.GetEmailCode(c, req.Email)
	if err != nil && !errors.Is(err, redis.Nil) {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Register] [GetEmailCode] [err] = ", err.Error())
		return
	}
	if code != req.Code {
		util.FailRespWithCodeAndZap(c, util.ReqDataError, "[Register] [GetEmailCode] [err] = ", "验证码错误")
		return
	}

	userId := dao.GetUserId(c, req.Email)
	if userId > 0 {
		util.FailRespWithCodeAndZap(c, util.ReqDataError, "[Register] [GetUserId] [err] = ", "用户已存在")
		return
	}

	if err := dao.CreateUser(c, req.Email, req.Username, req.Password); err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Register] [CreateUser] [err] = ", err.Error())
		return
	}

	userId = dao.GetUserId(c, req.Username)
	token, err := util.GenerateJWT(userId)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Register] [GenerateJWT] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, RegisterResp{Token: token})
}

func Info(c *gin.Context) {
	userId := util.GetUid(c)
	user, err := dao.GetUser(c, userId)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Info] [GetUser] [err] = ", err.Error())
		return
	}

	resp := &UserInfoResp{
		Email:     user.Email,
		UserName:  user.Username,
		StockList: make([]*UserInfoData, 0),
		FundList:  make([]*UserInfoData, 0),
	}

	rdb := connector.GetRedis().WithContext(c)
	eg, ctx := errgroup.WithContext(c)

	pipe := rdb.Pipeline()
	stockFollowCmd := pipe.SMembers(ctx, fmt.Sprintf(public.RedisKeyStockFollow, userId))
	fundFollowCmd := pipe.SMembers(ctx, fmt.Sprintf(public.RedisKeyFundFollow, userId))

	if _, err := pipe.Exec(ctx); err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Info] [Pipeline] [err] = ", err.Error())
		return
	}

	stockResult, _ := stockFollowCmd.Result()
	fundResult, _ := fundFollowCmd.Result()

	eg.Go(func() error {
		// 将股票关注列表转换为 int 列表
		stockIdList := cast.ToIntSlice(stockResult)
		if len(stockIdList) == 0 {
			return nil
		}

		// 获取股票信息
		stockInfos, err := dao.GetStockInfos(c, stockIdList)
		if err != nil {
			zap.S().Error("[Info] [GetStockInfos] [err] = ", err.Error())
			return err
		}

		// 组合 Redis Key
		predictKeys := make([]string, 0, len(stockInfos))
		for _, v := range stockInfos {
			predictKeys = append(predictKeys, fmt.Sprintf(public.RedisKeyStockPredict, v.Id))
		}
		pipe = rdb.Pipeline()
		predictCmd := pipe.MGet(ctx, predictKeys...)
		if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
			zap.S().Error("[Info] [Pipeline] [err] = ", err.Error())
			return err
		}
		predictResults, _ := predictCmd.Result()

		// 需要调用 Python 预测的股票
		var stocksToPredict []*model.StockInfo

		for i, v := range stockInfos {
			if predictResults[i] != nil {
				result, _ := rdb.Get(c, fmt.Sprintf(public.RedisKeyStockToday, v.TsCode)).Result()
				// 直接使用 Redis 预测值
				resp.StockList = append(resp.StockList, &UserInfoData{
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
		fundIdList := cast.ToIntSlice(fundResult)
		if len(fundIdList) == 0 {
			return nil
		}

		fundInfos, err := dao.GetFundInfos(c, fundIdList)
		if err != nil {
			zap.S().Error("[Info] [GetFundInfos] [err] = ", err.Error())
			return err
		}

		predictKeys := make([]string, 0, len(fundInfos))
		for _, v := range fundInfos {
			predictKeys = append(predictKeys, fmt.Sprintf(public.RedisKeyFundPredict, v.Id))
		}
		pipe = rdb.Pipeline()
		predictCmd := pipe.MGet(ctx, predictKeys...)
		if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
			zap.S().Error("[Info] [Pipeline] [err] = ", err.Error())
			return err
		}
		predictResults, _ := predictCmd.Result()

		var fundsToPredict []*model.FundInfo

		for i, v := range fundInfos {
			if predictResults[i] != nil {
				result, _ := rdb.Get(c, fmt.Sprintf(public.RedisKeyFundToday, v.TsCode)).Result()
				resp.FundList = append(resp.FundList, &UserInfoData{
					Id:      int(v.Id),
					Name:    v.Name,
					Val:     cast.ToFloat64(result),
					NextVal: cast.ToFloat64(predictResults[i]),
				})
			} else {
				fundsToPredict = append(fundsToPredict, v)
			}
		}

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
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "[Info] [Wait] [err] = ", err.Error())
		return
	}

	util.SuccessResp(c, resp)
}
