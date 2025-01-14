package dao

import (
	"context"
	"financia/public/db/connector"
	"financia/public/db/model"
)

func CreateFund(ctx context.Context, fund *model.FundInfo) error {
	return connector.GetDB().WithContext(ctx).Create(fund).Error
}

func InsertFund(ctx context.Context, data []*model.FundInfo) error {
	return connector.GetDB().WithContext(ctx).Create(data).Error
}
