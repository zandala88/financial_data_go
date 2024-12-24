package util

import (
	"crypto/md5"
	"financia/config"
	"fmt"
)

func GetMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s+config.Configs.App.Salt)))
}
