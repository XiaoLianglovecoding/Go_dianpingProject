package result

// Result 是整个项目统一返回给前端的 JSON 格式。
//
// Java 版黑马点评也是类似结构:
// { "success": true, "errorMsg": null, "data": ..., "total": null }
//
// 前端 axios 拦截器会检查 success：
// - success=true：继续渲染 data
// - success=false：弹出 errorMsg
type Result struct {
	// Success 表示接口是否成功。前端主要靠它判断是否继续处理数据。
	Success  bool   `json:"success"`
	ErrorMsg string `json:"errorMsg"` // ErrorMsg 是失败原因；成功时为空字符串。
	Data     any    `json:"data"`     // Data 是真正返回的数据，any 表示任意类型。
	Total    *int64 `json:"total"`    // Total 用于分页总数；指针为 nil 时 JSON 会显示 null。
}

// OK 返回一个没有 data 的成功结果。
func OK() Result {
	return Result{Success: true}
}

// OKWithData 返回一个带 data 的成功结果，最常用。
func OKWithData(data any) Result {
	return Result{Success: true, Data: data}
}

// OKWithList 返回列表数据和总数，后面做分页接口时会用到。
func OKWithList(data any, total int64) Result {
	return Result{Success: true, Data: data, Total: &total}
}

// Fail 返回失败结果。前端会把 message 弹出来。
func Fail(message string) Result {
	return Result{Success: false, ErrorMsg: message}
}
