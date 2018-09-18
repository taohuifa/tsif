package component

const (
	RESULT_FAIL    = 0 // 失败
	RESULT_SUCCESS = 1 // 成功
)

// 结果
type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 是否成功
func (this *Result) IsSuccess() bool {
	return this.Code > 0
}

func CreateResult(code int, msg string) *Result {
	return &Result{Code: code, Message: msg}
}

// 数据结果
type DataResult struct {
	Result
	Data interface{} `json:"data"`
}

func CreateDataResult(code int, msg string, data interface{}) *DataResult {
	return &DataResult{Result: Result{Code: code, Message: msg}, Data: data}
}
