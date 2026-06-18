package result

type Result struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"errorMsg"`
	Data     any    `json:"data"`
	Total    *int64 `json:"total"`
}

func OK() Result {
	return Result{Success: true}
}

func OKWithData(data any) Result {
	return Result{Success: true, Data: data}
}

func OKWithList(data any, total int64) Result {
	return Result{Success: true, Data: data, Total: &total}
}

func Fail(message string) Result {
	return Result{Success: false, ErrorMsg: message}
}
