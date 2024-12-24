package util

import (
	"errors"
	"financia/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

type MyClaims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(Id int64) (string, error) {
	c := MyClaims{
		Id, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.Configs.Auth.AccessExpire) * time.Second).Unix(), // 过期时间
			Issuer:    "zandala",                                                                            // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(config.Configs.Auth.AccessSecret))
}

func VerifyJWT(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.Configs.Auth.AccessSecret), nil
	})
	if err != nil {
		zap.S().Errorf("[VerifyJWT] 解析token失败: %v", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func GetUid(c *gin.Context) int64 {
	value, exists := c.Get("user_id")
	if !exists {
		zap.S().Info("[GetUid] ctx获取userId失败")
	}
	return cast.ToInt64(value)
}
