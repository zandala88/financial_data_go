package dao

import (
	"context"
	"financia/public/db/connector"
	"financia/public/db/model"
)

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

func GetFundList(ctx context.Context, search string, fundType, investType []string, page, pageSize int) ([]*model.FundInfo, int64, error) {
	var fundList []*model.FundInfo
	db := connector.GetDB().Model(&model.FundInfo{})
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

func GetFundInfo(ctx context.Context, id int) (*model.FundInfo, error) {
	var fundInfo model.FundInfo
	err := connector.GetDB().WithContext(ctx).Model(&model.FundInfo{}).
		Where("id = ?", id).First(&fundInfo).Error

	return &fundInfo, err
}

func CheckFundData(ctx context.Context, tsCode string) (bool, error) {
	var count int64
	err := connector.GetDB().WithContext(ctx).
		Raw("SELECT count(*) FROM t_fund_data WHERE f_ts_code = ?", tsCode).
		Scan(&count).Error

	return count > 0, err
}

func InsertFundData(ctx context.Context, data []*model.FundData) error {
	return connector.GetDB().WithContext(ctx).Create(data).Error
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
