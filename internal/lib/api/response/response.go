package response

// TODO: сделать корректные примеры в обход 3 структур для доки

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Error   string `json:"error" `
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewErrorResponse(errorMsg string, code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Error:   errorMsg,
		Code:    code,
		Message: message,
	}
}
