package user

type GetCodeReq struct {
	Email string `form:"email" binding:"required,email"`
}

type GetCodeResp struct {
	Code string `json:"code"`
}

type LoginReq struct {
	Email    string `form:"email"  binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type RegisterReq struct {
	Email    string `form:"email"  binding:"required,email"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Code     string `form:"code" binding:"required"`
}

type RegisterResp struct {
	Token string `json:"token"`
}

type UserInfoResp struct {
	Email     string          `json:"email"`
	UserName  string          `json:"username"`
	StockList []*UserInfoData `json:"stockList"`
	FundList  []*UserInfoData `json:"fundList"`
}

type UserInfoData struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Val     float64 `json:"val"`
	NextVal float64 `json:"nextVal"`
}

type TipResp struct {
	Exists    bool      `json:"exists"`
	StockRise TipSimple `json:"stockRise,omitempty"`
	StockFall TipSimple `json:"stockFall,omitempty"`
	FundRise  TipSimple `json:"fundRise,omitempty"`
	FundFall  TipSimple `json:"fundFall,omitempty"`
}

type TipSimple struct {
	Name  string  `json:"name"`
	Val   float64 `json:"val"`
	Scope float64 `json:"scope"`
}
