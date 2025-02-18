package dao

import (
	"context"
	"errors"
	"financia/public"
	"financia/public/db/connector"
	"financia/public/db/model"
	"financia/util"
	"fmt"
	"time"
)

// DistinctFundFields 获取基金字段
func DistinctFundFields(ctx context.Context) (map[string][]string, error) {
	var fundType []string
	var investType []string

	db := connector.GetDB()

	if err := db.Model(&model.FundInfo{}).WithContext(ctx).
		Distinct("f_fund_type").Pluck("f_fund_type", &fundType).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&model.FundInfo{}).WithContext(ctx).
		Distinct("f_invest_type").Where("f_invest_type != ''").Pluck("f_invest_type", &investType).Error; err != nil {
		return nil, err
	}

	// Combine results
	fields := map[string][]string{
		"fund_type":   fundType,
		"invest_type": investType,
	}

	return fields, nil
}

// GetFundList 获取基金列表
func GetFundList(ctx context.Context, search string, fundType, investType []string, page, pageSize int) ([]*model.FundInfo, int64, error) {
	var fundList []*model.FundInfo
	db := connector.GetDB().Model(&model.FundInfo{}).WithContext(ctx)
	if search != "" {
		db = db.Where("f_name like ?", "%"+search+"%").
			Or("f_management like ?", "%"+search+"%").
			Or("f_custodian like ?", "%"+search+"%").
			Or("f_trustee like ?", "%"+search+"%")
	}
	if len(fundType) > 0 {
		db = db.Where("f_fund_type in ?", fundType)
	}
	if len(investType) > 0 {
		db = db.Where("f_invest_type in ?", investType)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("f_flag desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&fundList).Error; err != nil {
		return nil, 0, err
	}

	return fundList, count, nil
}

// GetFundInfo 获取基金信息
func GetFundInfo(ctx context.Context, id int) (*model.FundInfo, error) {
	var fundInfo model.FundInfo
	err := connector.GetDB().WithContext(ctx).Model(&model.FundInfo{}).
		Where("id = ?", id).First(&fundInfo).Error

	return &fundInfo, err
}

func GetFundInfos(ctx context.Context, ids []int) ([]*model.FundInfo, error) {
	var fundInfos []*model.FundInfo
	err := connector.GetDB().WithContext(ctx).Model(&model.FundInfo{}).
		Where("id in ?", ids).Find(&fundInfos).Error

	return fundInfos, err
}

// CheckFundData 检查基金数据
func CheckFundData(ctx context.Context, tsCode string) (bool, error) {
	var count int64
	err := connector.GetDB().WithContext(ctx).
		Raw("SELECT count(*) FROM t_fund_data WHERE f_ts_code = ?", tsCode).
		Scan(&count).Error

	return count > 0, err
}

// InsertFundData 插入基金数据
func InsertFundData(ctx context.Context, data []*model.FundData) error {
	// 分批插入
	for i := 0; i < len(data); i += 1000 {
		end := i + 1000
		if end > len(data) {
			end = len(data)
		}
		if err := connector.GetDB().WithContext(ctx).Create(data[i:end]).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateFund(ctx context.Context, fund *model.FundInfo) error {
	return connector.GetDB().WithContext(ctx).Create(fund).Error
}

func InsertFund(ctx context.Context, data []*model.FundInfo) error {
	return connector.GetDB().WithContext(ctx).Create(data).Error
}

func UpdateFund(ctx context.Context, fund *model.FundInfo) error {
	return connector.GetDB().WithContext(ctx).Model(&model.FundInfo{}).
		Where("id = ?", fund.Id).Updates(fund).Error
}

// GetFundData 获取基金数据
func GetFundData(ctx context.Context, tsCode, start, end string) ([]*model.FundData, error) {
	var fundData []*model.FundData
	err := connector.GetDB().WithContext(ctx).
		Raw("SELECT * FROM t_fund_data WHERE f_ts_code = ? AND f_trade_date between ? AND ? order by f_trade_date", tsCode, start, end).
		Scan(&fundData).Error

	return fundData, err
}

func GetFundDataLimit30(ctx context.Context, tsCode string) ([]*model.FundData, error) {
	fundData := make([]*model.FundData, 0)
	err := connector.GetDB().WithContext(ctx).
		Raw("SELECT * FROM t_fund_data WHERE f_ts_code = ? order by f_trade_date desc limit 31", tsCode).
		Scan(&fundData).Error

	if len(fundData) == 0 {
		return nil, errors.New("no stock data")
	}

	rdb := connector.GetRedis().WithContext(ctx)
	rdb.Set(ctx, fmt.Sprintf(public.RedisKeyFundToday, tsCode), fundData[0].Close, time.Duration(util.SecondsUntilMidnight())*time.Second)

	return fundData, err
}
