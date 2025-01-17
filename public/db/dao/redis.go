package dao

import (
	"context"
	"encoding/json"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/model"
	"fmt"
	"time"
)

func SetEmailCode(ctx context.Context, email, code string) error {
	return connector.GetRedis().Set(ctx, email, code, 60*time.Second).Err()
}

func GetEmailCode(ctx context.Context, email string) (string, error) {
	return connector.GetRedis().Get(ctx, email).Result()
}

func StoreYearData(ctx context.Context, year string, values []float64) error {
	// 将数据结构转换为 JSON 字符串
	yearData := model.YearData{
		Year: year,
		Data: values,
	}
	jsonData, err := json.Marshal(yearData)
	if err != nil {
		return fmt.Errorf("could not marshal data: %v", err)
	}

	// 将 JSON 字符串存入 Redis List
	err = connector.GetRedis().LPush(ctx, public.RedisKeyFundSalesRatio, jsonData).Err()
	if err != nil {
		return fmt.Errorf("could not push to Redis list: %v", err)
	}
	return nil
}

func GetAllYearData(ctx context.Context) ([]model.YearData, error) {
	// 获取 Redis 中指定的 List 中所有的数据
	jsonDataList, err := connector.GetRedis().LRange(ctx, public.RedisKeyFundSalesRatio, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("could not get all data from Redis list: %v", err)
	}

	// 反序列化所有 JSON 数据
	var allYearData []model.YearData
	for _, jsonData := range jsonDataList {
		var yearData model.YearData
		err := json.Unmarshal([]byte(jsonData), &yearData)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal data: %v", err)
		}
		allYearData = append(allYearData, yearData)
	}

	return allYearData, nil
}
