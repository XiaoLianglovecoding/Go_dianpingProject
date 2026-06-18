package dto

// ScrollResult 是滚动分页结果，后面做“关注流/Feed流”时会用。
//
// 普通分页用 page/current；Feed 流常用 minTime + offset 来继续向下翻。
type ScrollResult struct {
	List    any   `json:"list"` // 当前页数据。
	MinTime int64 `json:"minTime"`
	Offset  int   `json:"offset"` // 同一时间戳下的偏移量。
}
