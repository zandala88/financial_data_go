package tushare

import (
	"encoding/json"
	"financia/config"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const (
	url        = "http://api.tushare.pro"
	timeLayout = "20060102"
)

var token = config.Configs.TuShare.Token

type TuShareReq struct {
	ApiName string      `json:"api_name"`
	Token   string      `json:"token"`
	Params  interface{} `json:"params"`
	Fields  string      `json:"fields"`
}

type TuShareResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func tuSharePost(apiName string, data interface{}, fields string) any {
	client := resty.New()
	tuShareResp := &TuShareResp{}

	zap.S().Debugf("[tuSharePost] [apiName] = %s", apiName)
	zap.S().Debugf("[tuSharePost] [data] = %#v", data)

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&TuShareReq{
			ApiName: apiName,
			Token:   token,
			Params:  data,
			Fields:  fields,
		}).
		SetResult(&tuShareResp).
		Post(url)
	if err != nil {
		zap.S().Errorf("[tuSharePost] [err] = %s", err.Error())
		return nil
	}

	return tuShareResp.Data
}

func marshalResp(r any, resp *DailyResp) error {
	marshal, err := json.Marshal(r.(map[string]interface{}))
	if err != nil {
		zap.S().Errorf("[FutWeeklyDetail] [json.Marshal] [err] = %s", err.Error())
		return err
	}

	if err := json.Unmarshal(marshal, &resp); err != nil {
		zap.S().Errorf("[FutWeeklyDetail] [json.Unmarshal] [err] = %s", err.Error())
		return err
	}
	return nil
}
