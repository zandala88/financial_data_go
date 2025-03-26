package spark

type SparkReq struct {
	Model    string         `json:"model"`
	User     string         `json:"user"`
	Messages []SparkMessage `json:"messages"`
	Stream   bool           `json:"stream"`
}

type SparkMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SparkResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Choices []struct {
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		Index int `json:"index"`
	} `json:"choices"`
}

type AnalyzeData struct {
	Data    []float64 `json:"data"`
	SMA     []float64 `json:"sma"`
	EMA     []float64 `json:"ema"`
	WMA     []float64 `json:"wma"`
	MACD    []float64 `json:"macd"`
	RSI     []float64 `json:"rsi"`
	Predict float64   `json:"predict"`
}
