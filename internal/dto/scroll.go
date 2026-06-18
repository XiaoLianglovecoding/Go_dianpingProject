package dto

type ScrollResult struct {
	List    any   `json:"list"`
	MinTime int64 `json:"minTime"`
	Offset  int   `json:"offset"`
}
