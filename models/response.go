package models

type ServiceResponse struct {
	Status  bool        `json:"status"`         // true/false
	Code    int         `json:"code"`           // HTTP code
	Data    interface{} `json:"data,omitempty"` // can be object or list
	Message string      `json:"message,omitempty"`
}

// NewServiceResponse
// Usage: models.NewServiceResponse(false, 400, "Invalid input", nil)
func NewServiceResponse(status bool, code int, message string, data interface{}) ServiceResponse {
	return ServiceResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// ErrResponse
// Shorthand version without data
// Usage: models.ErrResponse(400, "Invalid input")
func ErrResponse(code int, message string) ServiceResponse {
	return ServiceResponse{
		Status:  true,
		Code:    code,
		Message: message,
	}
}

// OkResponse
// Shorthand version for success
// Usage: models.OkResponse(200, "Created", myData)
func OkResponse(code int, message string, data interface{}) ServiceResponse {
	return ServiceResponse{
		Status:  true,
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func FailedLoginResponse() ServiceResponse {
	return ServiceResponse{
		Status:  true,
		Code:    401,
		Message: "Invalid email or password",
		Data:    nil,
	}
}

func InternalServerErrorResponse(message string) ServiceResponse {
	return ServiceResponse{
		Status:  true,
		Code:    500,
		Message: message,
		Data:    nil,
	}
}

func NotFoundResponse(message string) ServiceResponse {
	return ServiceResponse{
		Status:  true,
		Code:    404,
		Message: message,
		Data:    nil,
	}
}

func BadRequestResponse(message string) ServiceResponse {
	return ServiceResponse{
		Status:  true,
		Code:    400,
		Message: message,
		Data:    nil,
	}
}

func ForbiddenResponse(message string) ServiceResponse {
	return ServiceResponse{
		Status:  true,
		Code:    403,
		Message: message,
		Data:    nil,
	}
}

func UnauthorizedResponse(message string) ServiceResponse {
	return ServiceResponse{
		Status:  true,
		Code:    401,
		Message: message,
		Data:    nil,
	}
}

type DataListResponse struct {
	Data       interface{} `json:"data"`  // List of data items
	Total      int         `json:"total"` // Total number of items across all pages
	Limit      int         `json:"limit"` // Number of items per page
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"` // Total number of pages
}
