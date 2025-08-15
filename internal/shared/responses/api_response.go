package shared

type ApiResponse[T any] struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Content *T          `json:"content,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SuccessResponse[T any](message string, content *T) ApiResponse[T] {
	return ApiResponse[T]{
		Success: true,
		Message: message,
		Content: content,
		Error:   nil,
	}
}

func ErrorResponse(message string, err interface{}) ApiResponse[any] {
	return ApiResponse[any]{
		Success: false,
		Message: message,
		Content: nil,
		Error:   err,
	}
}
