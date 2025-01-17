package dao

import (
	"context"
	"financia/public/db/connector"
	"financia/public/db/model"
	"fmt"
	"gorm.io/gorm"
)

func DistinctStockFields(ctx context.Context) (map[string][]string, error) {
	var isHs []string
	var exchange []string
	var market []string

	db := connector.GetDB()

	// Query distinct f_is_hs
	if err := db.Model(&model.StockInfo{}).WithContext(ctx).
		Distinct("f_is_hs").Pluck("f_is_hs", &isHs).Error; err != nil {
		return nil, err
	}

	// Query distinct f_exchange
	if err := db.Model(&model.StockInfo{}).WithContext(ctx).
		Distinct("f_exchange").Pluck("f_exchange", &exchange).Error; err != nil {
		return nil, err
	}

	// Query distinct f_market
	if err := db.Model(&model.StockInfo{}).WithContext(ctx).
		Distinct("f_market").Pluck("f_market", &market).Error; err != nil {
		return nil, err
	}

	// Combine results
	fields := map[string][]string{
		"is_hs":    isHs,
		"exchange": exchange,
		"market":   market,
	}

	return fields, nil
}

func CountStockFields(ctx context.Context) (map[string]map[string]int, error) {
	queryByField := func(db *gorm.DB, field string) map[string]int {
		type StockInfo struct {
			Value string `gorm:"column:value"` // 动态字段的通用名称
			Count int    `gorm:"column:cnt"`
		}

		var results []StockInfo

		db.Model(&model.StockInfo{}).
			Select(fmt.Sprintf("%s AS value, COUNT(1) as cnt", field)). // 使用 AS 设置别名
			Group(field).
			Scan(&results)

		mp := make(map[string]int)

		for _, result := range results {
			mp[result.Value] = result.Count
		}
		return mp
	}

	db := connector.GetDB()
	mp1 := queryByField(db, "f_is_hs")
	mp2 := queryByField(db, "f_exchange")
	mp3 := queryByField(db, "f_market")

	fields := map[string]map[string]int{
		"is_hs":    mp1,
		"exchange": mp2,
		"market":   mp3,
	}

	return fields, nil
}

func GetStockInfo(ctx context.Context, id int) (*model.StockInfo, error) {
	var stockInfo model.StockInfo
	err := connector.GetDB().WithContext(ctx).Model(&model.StockInfo{}).
		Where("f_id = ?", id).First(&stockInfo).Error

	return &stockInfo, err
}

func GetStockList(ctx context.Context, search string, isHs, exchange, market []string, page, pageSize int) ([]*model.StockInfo, int64, error) {
	var stockList []*model.StockInfo
	db := connector.GetDB().Model(&model.StockInfo{})
	if search != "" {
		db = db.Where("f_name like ?", "%"+search+"%").
			Or("f_fullname like ?", "%"+search+"%")
	}
	if len(isHs) > 0 {
		db = db.Where("f_is_hs in ?", isHs)
	}
	if len(exchange) > 0 {
		db = db.Where("f_exchange in ?", exchange)
	}
	if len(market) > 0 {
		db = db.Where("f_market in ?", market)
	}

	var count int64
	err := db.WithContext(ctx).Count(&count).Scopes(Paginate(page, pageSize)).Order("f_name DESC").Find(&stockList).Error

	return stockList, count, err
}

func GetStockData(ctx context.Context, tsCode, start, end string) ([]*model.StockData, error) {
	var stockData []*model.StockData
	err := connector.GetDB().WithContext(ctx).
		Raw("SELECT * FROM t_stock_data WHERE f_ts_code = ? AND f_trade_date between ? AND ? order by f_trade_date", tsCode, start, end).
		Scan(&stockData).Error

	return stockData, err
}

func CheckStockData(ctx context.Context, tsCode string) (bool, error) {
	var count int64
	err := connector.GetDB().WithContext(ctx).
		Raw("SELECT count(*) FROM t_stock_data WHERE f_ts_code = ?", tsCode).
		Scan(&count).Error

	return count > 0, err
}

func CreateStockData(ctx context.Context, data *model.StockData) error {
	return connector.GetDB().WithContext(ctx).Create(data).Error
}

func InsertStockData(ctx context.Context, data []*model.StockData) error {
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

// SameMarketList 查询同市场股票列表
func SameMarketList(ctx context.Context, market string, id int) ([]*model.StockInfo, error) {
	var stockList []*model.StockInfo
	err := connector.GetDB().WithContext(ctx).
		Model(&model.StockInfo{}).Where("f_market = ?", market).Where("f_id <> ?", id).Find(&stockList).Error
	return stockList, err
}

func SameIndustryList(ctx context.Context, industry string, id int) ([]*model.StockInfo, error) {
	var stockList []*model.StockInfo
	err := connector.GetDB().WithContext(ctx).
		Model(&model.StockInfo{}).Where("f_industry = ?", industry).Where("f_id <> ?", id).Find(&stockList).Error
	return stockList, err
}
