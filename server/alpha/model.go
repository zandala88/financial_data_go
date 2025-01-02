package alpha

type StockMetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

type StockDailyData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type StockTimeSeriesDaily map[string]StockDailyData

type StockResp struct {
	MetaData        StockMetaData        `json:"Meta Data"`
	TimeSeriesDaily StockTimeSeriesDaily `json:"Time Series (Daily)"`
}

// ----------------------------

type CurrencyMetaData struct {
	Information   string `json:"1. Information"`
	FromSymbol    string `json:"2. From Symbol"`
	ToSymbol      string `json:"3. To Symbol"`
	OutputSize    string `json:"4. Output Size"`
	LastRefreshed string `json:"5. Last Refreshed"`
	TimeZone      string `json:"6. Time Zone"`
}

type CurrencyDailyData struct {
	Open  string `json:"1. open"`
	High  string `json:"2. high"`
	Low   string `json:"3. low"`
	Close string `json:"4. close"`
}

type CurrencyTimeSeriesDaily map[string]CurrencyDailyData

type CurrencyResp struct {
	MetaData        CurrencyMetaData        `json:"Meta Data"`
	TimeSeriesDaily CurrencyTimeSeriesDaily `json:"Time Series FX (Daily)"`
}

// ----------------------------

type CryptoMetaData struct {
	Information   string `json:"1. Information"`
	DigitalCode   string `json:"2. Digital Currency Code"`
	DigitalSymbol string `json:"3. Digital Currency Name"`
	MarketCode    string `json:"4. Market Code"`
	MarketSymbol  string `json:"5. Market Name"`
	LastRefreshed string `json:"6. Last Refreshed"`
	TimeZone      string `json:"7. Time Zone"`
}

type CryptoDailyData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type CryptoTimeSeriesDaily map[string]CryptoDailyData

type CryptoResp struct {
	MetaData        CryptoMetaData        `json:"Meta Data"`
	TimeSeriesDaily CryptoTimeSeriesDaily `json:"Time Series (Digital Currency Daily)"`
}

// ----------------------------

type LstmResp struct {
	Name        string    `json:"name"`
	Predictions []float64 `json:"predictions"`
}
