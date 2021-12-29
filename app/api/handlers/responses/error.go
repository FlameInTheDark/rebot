package responses

type ResponseError struct {
	Message string `json:"message"`
}

func Error(err error) ResponseError {
	return ResponseError{Message: err.Error()}
}
