package util

// SMA 简单移动平均线
func SMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil // 数据不足
	}

	sma := make([]float64, len(prices)-period+1)
	sum := 0.0

	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	sma[0] = sum / float64(period)

	for i := period; i < len(prices); i++ {
		sum += prices[i] - prices[i-period]
		sma[i-period+1] = sum / float64(period)
	}

	return sma
}

// EMA 指数移动平均线
func EMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil // 数据不足
	}

	ema := make([]float64, len(prices))
	alpha := 2.0 / (float64(period) + 1.0)

	// 初始化：第一个 EMA 使用 SMA 计算
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	// 计算 EMA
	for i := period; i < len(prices); i++ {
		ema[i] = alpha*prices[i] + (1-alpha)*ema[i-1]
	}

	return ema[period-1:] // 返回有效数据部分
}

// WMA 加权移动平均线
func WMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil // 数据不足
	}

	wma := make([]float64, len(prices)-period+1)
	weightSum := float64((period * (period + 1)) / 2)

	for i := period - 1; i < len(prices); i++ {
		weightedSum := 0.0
		for j := 0; j < period; j++ {
			weightedSum += prices[i-j] * float64(period-j)
		}
		wma[i-period+1] = weightedSum / weightSum
	}

	return wma
}

// MACD 指数平滑异同平均线
func MACD(prices []float64, shortPeriod, longPeriod, signalPeriod int) ([]float64, []float64, []float64) {
	if len(prices) < longPeriod {
		return nil, nil, nil // 数据不足
	}

	emaShort := EMA(prices, shortPeriod)
	emaLong := EMA(prices, longPeriod)

	// 保证 emaShort 和 emaLong 长度一致（取相同的末尾部分）
	minLen := min(len(emaShort), len(emaLong))
	emaShort = emaShort[len(emaShort)-minLen:]
	emaLong = emaLong[len(emaLong)-minLen:]

	macdLine := make([]float64, minLen)
	for i := 0; i < minLen; i++ {
		macdLine[i] = emaShort[i] - emaLong[i]
	}

	signalLine := EMA(macdLine, signalPeriod)
	if len(signalLine) == 0 {
		return macdLine, nil, nil
	}

	// 保证 histogram 计算不会越界
	histogram := make([]float64, len(signalLine))
	for i := 0; i < len(signalLine); i++ {
		histogram[i] = macdLine[len(macdLine)-len(signalLine)+i] - signalLine[i]
	}

	return macdLine[len(macdLine)-len(signalLine):], signalLine, histogram
}

// RSI 相对强弱指标
func RSI(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	rsi := make([]float64, len(prices)-period+1)
	gains, losses := 0.0, 0.0

	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)
	rs := avgGain / avgLoss
	rsi[0] = 100 - (100 / (1 + rs))

	for i := period; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			avgGain = ((avgGain * float64(period-1)) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = ((avgLoss * float64(period-1)) - change) / float64(period)
		}

		rs = avgGain / avgLoss
		rsi[i-period+1] = 100 - (100 / (1 + rs))
	}

	return rsi
}
