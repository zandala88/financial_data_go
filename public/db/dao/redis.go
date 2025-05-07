package dao

import (
	"context"
	"financia/public"
	"financia/public/db/connector"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

func SetEmailCode(ctx context.Context, email, code string) error {
	return connector.GetRedis().Set(ctx, email, code, 60*time.Second).Err()
}

func GetEmailCode(ctx context.Context, email string) (string, error) {
	return connector.GetRedis().Get(ctx, email).Result()
}

func GetFollowList(c context.Context, userId int64) ([]int, []int, error) {
	rdb := connector.GetRedis().WithContext(c)

	pipe := rdb.Pipeline()
	stockFollowCmd := pipe.SMembers(c, fmt.Sprintf(public.RedisKeyStockFollow, userId))
	fundFollowCmd := pipe.SMembers(c, fmt.Sprintf(public.RedisKeyFundFollow, userId))

	if _, err := pipe.Exec(c); err != nil {
		zap.S().Error("[Info] [Pipeline] [err] = ", err.Error())
		return nil, nil, err
	}

	stockResult, _ := stockFollowCmd.Result()
	fundResult, _ := fundFollowCmd.Result()

	return cast.ToIntSlice(stockResult), cast.ToIntSlice(fundResult), nil
}
