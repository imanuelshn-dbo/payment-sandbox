package apperror

type AppError struct {
	Code    int               `json:"-"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// helper constructors
func BadRequest(msg string) *AppError {
	return &AppError{Code: 400, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: 401, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Code: 403, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: 404, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: 409, Message: msg}
}

func Internal(msg string) *AppError {
	return &AppError{Code: 500, Message: msg}
}

func Validation(errs map[string]string) *AppError {
	return &AppError{
		Code:    400,
		Message: "validation error",
		Errors:  errs,
	}
}
