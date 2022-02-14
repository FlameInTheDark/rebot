package responses

// ResponseError is an API error message
type ResponseError struct {
	Message string `json:"message"`
}

// Error creates new ResponseError with error message
func Error(err error) ResponseError {
	return ResponseError{Message: err.Error()}
}
