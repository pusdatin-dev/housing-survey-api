package shared

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New()
}

// CustomValidate handles common error formatting
func CustomValidate(input interface{}, customMessages map[string]string) error {
	err := Validator.Struct(input)
	if err == nil {
		return nil
	}

	var invalidValidationError *validator.InvalidValidationError
	if ok := errors.As(err, &invalidValidationError); ok {
		return fmt.Errorf("invalid input structure")
	}

	var messages []string
	for _, fieldErr := range err.(validator.ValidationErrors) {
		key := fmt.Sprintf("%s.%s", fieldErr.Field(), fieldErr.Tag())
		if msg, ok := customMessages[key]; ok {
			messages = append(messages, msg)
		} else {
			messages = append(messages, fmt.Sprintf("%s is invalid (%s)", fieldErr.Field(), fieldErr.Tag()))
		}
	}

	return fmt.Errorf(strings.Join(messages, "; "))
}
