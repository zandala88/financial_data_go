package dao

import (
	"context"
	"financia/public/db/connector"
	"time"
)

func SetEmailCode(ctx context.Context, email, code string) error {
	return connector.GetRedis().Set(ctx, email, code, 60*time.Second).Err()
}

func GetEmailCode(ctx context.Context, email string) (string, error) {
	return connector.GetRedis().Get(ctx, email).Result()
}
