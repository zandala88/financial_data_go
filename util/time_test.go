package util

import "testing"

func Test_TimeParse(t *testing.T) {
	t.Log(ConvertDateStrToTime("20210801", "20060102"))
}
