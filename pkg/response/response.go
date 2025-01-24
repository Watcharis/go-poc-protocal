package response

const (
	STATUS_SUCCESS = "success"
	STATUS_ERROR   = "error"
)

type CommonResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func SetCommonResponse(status string, code int) CommonResponse {
	return CommonResponse{
		Status: status,
		Code:   code,
	}
}
