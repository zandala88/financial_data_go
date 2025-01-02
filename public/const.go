package public

// 1 股票 2 外汇
const (
	StockType = iota + 1
	CurrencyType
	CryptoType
)

// 1.LSTM 2.ARIMA
const (
	LstmType = iota + 1
	ArimaType
)
