package validation

import (
	"gopkg.in/go-playground/validator.v9"
)

type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

// ValidateStruct validates struct date and returns errors
func ValidateStruct(v interface{}) []*ErrorResponse {
	var errs []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		if validations, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validations {
				var element ErrorResponse
				element.FailedField = err.StructNamespace()
				element.Tag = err.Tag()
				element.Value = err.Param()
				errs = append(errs, &element)
			}
		} else {
			errs = append(errs, &ErrorResponse{
				FailedField: "validation check error",
				Tag:         "validation check error",
				Value:       "validation check error",
			})
		}
	}
	return errs
}
