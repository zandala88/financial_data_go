package company

type DetailCompanyReq struct {
	Id int `form:"id" binding:"required"`
}

type DetailCompanyResp struct {
	ComName       string  `json:"comName"`
	ComId         string  `json:"comId"`
	Chairman      string  `json:"chairman"`
	Manager       string  `json:"manager"`
	Secretary     string  `json:"secretary"`
	RegCapital    float64 `json:"regCapital"`
	Province      string  `json:"province"`
	City          string  `json:"city"`
	Employees     int     `json:"employees"`
	Introduction  string  `json:"introduction"`
	BusinessScope string  `json:"businessScope"`
	MainBusiness  string  `json:"mainBusiness"`
}

type ListCompanyReq struct {
	Search   string   `form:"search"`
	Province []string `form:"province"`
	Page     int      `form:"page" binding:"required"`
	PageSize int      `form:"pageSize" binding:"required"`
}

type ListCompanyResp struct {
	List         []*ListCompanySimple `json:"list"`
	TotalPageNum int                  `json:"totalPageNum"`
	HasMore      bool                 `json:"hasMore"`
	Count        int64                `json:"count"`
}

type ListCompanySimple struct {
	Id         int     `json:"id"`
	ComName    string  `json:"comName"`
	ComId      string  `json:"comId"`
	Chairman   string  `json:"chairman"`
	Manager    string  `json:"manager"`
	Secretary  string  `json:"secretary"`
	RegCapital float64 `json:"regCapital"`
	Province   string  `json:"province"`
	City       string  `json:"city"`
	Employees  int     `json:"employees"`
}

type QueryCompanyResp struct {
	List []string `json:"list"`
}
