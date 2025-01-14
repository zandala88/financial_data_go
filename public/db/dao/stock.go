package dao

import (
	"context"
	"financia/public/db/connector"
	"financia/public/db/model"
)

func DistinctFields(ctx context.Context) (map[string][]string, error) {
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
	return connector.GetDB().WithContext(ctx).Create(&data).Error
}
