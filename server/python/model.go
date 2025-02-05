package python

type PythonPredictReq struct {
	Data []*PythonPredictReqSimple `json:"data"`
}

type PythonPredictReqSimple struct {
	Date   string  `json:"date"`
	CoIMF1 float64 `json:"Co-IMF1"`
	CoIMF2 float64 `json:"Co-IMF2"`
	CoIMF3 float64 `json:"Co-IMF3"`
	CoIMF4 int64   `json:"Co-IMF4"`
	Target float64 `json:"Target"`
}

type PythonPredictResp struct {
	Code int                   `json:"code"`
	Data PythonPredictRespData `json:"data"`
}

type PythonPredictRespData struct {
	Val float64 `json:"val"`
}
