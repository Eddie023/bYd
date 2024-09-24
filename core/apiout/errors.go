package apiout

// ErrorResponse defines a standard application error.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CustomError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface.
func (e *CustomError) Error() string {
	return e.Message
}
