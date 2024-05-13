package csval

import (
	"net/mail"
	"net/url"
)

type ValidationResult struct {
	Pass     bool
	Messages []string
}

// NewSuccessValidationResult return a new validation result
func NewSuccessValidationResult() ValidationResult {
	return ValidationResult{Pass: true}
}

// NewFailingValidationResult return a new validation result
func NewFailingValidationResult(msg ...string) ValidationResult {
	return ValidationResult{Pass: false, Messages: msg}
}

// Merge aggregate the outcome of 2 validation results
func (v *ValidationResult) Append(result ValidationResult) {
	if result.Pass {
		return
	}

	v.Pass = false
	v.Messages = append(v.Messages, result.Messages...)
}

// Exists return a passing result if the string has a valid value
func IsNotEmpty(input string) ValidationResult {
	if len(input) > 0 {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult("input cannot be empty")
}

// IsEmail return success if the string is a valid email format
func IsEmail(input string) ValidationResult {
	_, err := mail.ParseAddress(input)
	if err == nil {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult("email is not valid")
}

// IsValidWebAddress return success if the string is a valid email format
func IsValidWebAddress(input string) ValidationResult {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult("web address is not valid")
}

// IsGreaterThan  return success if the value of the number is less than the maximum
func IsGreaterThan(input int, target int) ValidationResult {
	if input > target {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult("number is less than allowed")
}

// IsLessThan  return success if the value of the number is greater than the minimum
func IsLessThan(input int, target int) ValidationResult {
	if input < target {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult("number is greater than than allowed")
}

// IsLengthGreaterThan return success if the length of the string is greater than the minimum
func IsLengthGreaterThan(input string, target int) ValidationResult {
	if len(input) > target {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult("number is less than allowed")
}

// IsLengthLessThan return success if the length of the string is less than the maximum
func IsLengthLessThan(input string, target int) ValidationResult {
	if len(input) < target {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult("number is greater than than allowed")
}
