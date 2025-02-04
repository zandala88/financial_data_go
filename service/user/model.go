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
