package dao

import (
	"context"
	"financia/public/db/connector"
	"financia/public/db/model"
)

// ProvinceDis 获取省份分布
// 筛选参数
func ProvinceDis(ctx context.Context) ([]string, error) {
	var provinces []string
	err := connector.GetDB().WithContext(ctx).
		Model(&model.CompanyInfo{}).
		Distinct("f_province").
		Pluck("f_province", &provinces).Error
	return provinces, err
}

// GetCompanyList 获取公司列表
func GetCompanyList(ctx context.Context, search string, province []string, page, pageSize int) ([]*model.CompanyInfo, int64, error) {
	var companyList []*model.CompanyInfo
	db := connector.GetDB().Model(&model.CompanyInfo{})
	if search != "" {
		db = db.Where("f_com_name like ?", "%"+search+"%").
			Or("f_chairman like ?", "%"+search+"%").
			Or("f_manager like ?", "%"+search+"%").
			Or("f_secretary like ?", "%"+search+"%")
	}
	if len(province) > 0 {
		db = db.Where("f_province in ?", province)
	}

	var count int64
	err := db.WithContext(ctx).Count(&count).Scopes(Paginate(page, pageSize)).Order("f_com_name DESC,f_com_id DESC").Find(&companyList).Error

	return companyList, count, err
}

// GetCompany 获取公司信息
func GetCompany(ctx context.Context, id int) (*model.CompanyInfo, error) {
	var company *model.CompanyInfo
	err := connector.GetDB().WithContext(ctx).Where("f_id = ?", id).First(&company).Error
	return company, err
}
