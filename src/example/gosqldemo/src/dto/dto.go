package dto

// 公共返回体
type DataResponse struct {
	Success bool
	Message string
	Data    interface{}
}

func NewDataResponse() DataResponse {
	dr := DataResponse{}
	dr.Success = dr.Message == ""
	return dr
}

func (dr DataResponse) SetMessage(msg string) {
	dr.Message = msg
	dr.Success = false
}
