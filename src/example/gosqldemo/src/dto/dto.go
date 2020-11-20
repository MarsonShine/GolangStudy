package dto

// 公共返回体
type DataResponse struct {
	Success bool
	Message string
	Data    interface{}
}
