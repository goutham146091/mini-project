package view

// Response template response structure returned by api
type Response struct {
	Status    int      `json:"status"`
	Message   string      `json:"message"`
	ErrorCode int         `json:"error_code"`
	Payload   interface{} `json:"payload"`
}

// GetFailureResponse returns a formatted error response in case on any failure
func GetFailureResponse(m string, c int, payload interface{}) *Response {

	return &Response{
		Status:    4000,
		Message:   m,
		ErrorCode: c,
		Payload:    payload,
	}
}

// GetSuccessResponse returns a formatted success response
func GetSuccessResponse(m string, payload interface{}) *Response {

	return &Response{
		Status:    2000,
		Message:   m,
		ErrorCode: 0,
		Payload:    payload,
	}
}

// GetSuccessResponseWithErrorCode returns a formatted success response
func GetSuccessResponseWithErrorCode(m string, errorCode int, payload interface{}) *Response {

	return &Response{
		Status:    2000,
		Message:   m,
		ErrorCode: errorCode,
		Payload:    payload,
	}
}