package spark

import (
	"bufio"
	"encoding/json"
	"financia/config"
	"financia/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"io"
	"strings"
)

var (
	url     = "https://spark-api-open.xf-yun.com/v1/chat/completions"
	token   = "Bearer %s"
	model   = "4.0Ultra"
	predMsg = "你是一个专业的金融量化分析师。我会给你股票或者基金一定时间内的收盘价、简单移动平均线、指数移动平均线、加权移动平均线、" +
		"指数平滑异同平均线、相对强弱指标，请你分析这些数据并用中文回答我。"
)

func SendSparkHttp(c *gin.Context, arr []float64, userId string) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许跨域

	_, _, macd := util.MACD(arr, 5, 10, 5)
	msg := &AnalyzeData{
		Data: arr,
		SMA:  util.SMA(arr, 5),
		EMA:  util.EMA(arr, 5),
		WMA:  util.WMA(arr, 5),
		MACD: macd,
		RSI:  util.RSI(arr, 5),
	}
	msgStr, _ := json.Marshal(msg)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf(token, config.Configs.Spark.Password)).
		SetBody(&SparkReq{
			Model: model,
			User:  userId,
			Messages: []SparkMessage{
				{
					Role:    "system",
					Content: predMsg,
				},
				{
					Role:    "user",
					Content: string(msgStr),
				},
			},
			Stream: true,
		}).
		SetDoNotParseResponse(true).
		Post(url)
	if err != nil {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "请求失败", err.Error())
		return
	}
	defer resp.RawBody().Close()

	reader := bufio.NewScanner(resp.RawBody())
	for reader.Scan() {
		line := reader.Text()
		if strings.HasPrefix(line, "data:") {
			// 去掉 "data:" 前缀并解析 JSON
			jsonStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			var sparkResp SparkResp
			if err := json.Unmarshal([]byte(jsonStr), &sparkResp); err != nil {
				continue
			}

			// 逐条发送 SSE 数据
			for _, choice := range sparkResp.Choices {
				data := choice.Delta.Content
				fmt.Fprintf(c.Writer, "%s", data)
				c.Writer.Flush() // 立即推送数据
			}
		}
	}

	if err := reader.Err(); err != nil && err != io.EOF {
		util.FailRespWithCodeAndZap(c, util.InternalServerError, "流式处理失败", err.Error())
	}

}
