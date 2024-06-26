package validate

import (
	"errors"
	"fmt"

	"net"
	"net/mail"
	"net/url"

	"regexp"
)

// ValidationResult an object that holds the aggregated outcome of validation routines
type ValidationResult struct {
	Pass     bool
	Messages []ValidationMessage
}

// ValidationMessage a struct that holds a message and the field that it relates to
type ValidationMessage struct {
	Message string
	Field   string
}

// NewSuccessValidationResult return a new validation result
func NewSuccessValidationResult() ValidationResult {
	return ValidationResult{Pass: true}
}

// NewFailingValidationResult return a new validation result
func NewFailingValidationResult(msg ...ValidationMessage) ValidationResult {
	return ValidationResult{Pass: false, Messages: msg}
}

// Append aggregate the outcome of 2 validation results
func (v *ValidationResult) Append(result ValidationResult) {
	if result.Pass {
		return
	}

	v.Pass = false
	v.Messages = append(v.Messages, result.Messages...)
}

// Error if errors are available, return an aggregated single error
func (v *ValidationResult) Error(context string) error {
	if v.Pass {
		return nil
	}

	out := []error{fmt.Errorf(context)}
	for _, msg := range v.Messages {
		out = append(out, fmt.Errorf("%s: %s", msg.Field, msg.Message))
	}

	return errors.Join(out...)
}

// NewValidationMessage return a validation message
func NewValidationMessage(field string, msg string) ValidationMessage {
	return ValidationMessage{Message: msg, Field: field}
}

// IsNotEmpty return a passing result if the string has a valid value
func IsNotEmpty(field string, input string) ValidationResult {
	if len(input) > 0 {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, "input cannot be empty"))
}

// IsEmail return success if the string is a valid email format
func IsEmail(field string, input string) ValidationResult {
	_, err := mail.ParseAddress(input)
	if err == nil {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, "email is not valid"))
}

// IsValidWebAddress return success if the string is a valid email format
func IsValidWebAddress(field string, input string) ValidationResult {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, "web address is not valid"))
}

// IsIP return success if the string is a valid IP format
func IsIP(field string, input string) ValidationResult {
	ip := net.ParseIP(input)
	if ip != nil {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, "IP is not valid"))
}

// IsGreaterThan  return success if the value of the number is less than the maximum
func IsGreaterThan(field string, input int, target int) ValidationResult {
	if input > target {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, "number is less than allowed"))
}

// IsLessThan  return success if the value of the number is greater than the minimum
func IsLessThan(field string, input int, target int) ValidationResult {
	if input < target {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, "number is greater than than allowed"))
}

// IsLengthGreaterThan return success if the length of the string is greater than the minimum
func IsLengthGreaterThan(field string, input string, target int) ValidationResult {
	if len(input) > target {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, fmt.Sprintf("string length is less than %v", target)))
}

// IsLengthLessThan return success if the length of the string is less than the maximum
func IsLengthLessThan(field string, input string, target int) ValidationResult {
	if len(input) < target {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, fmt.Sprintf("string length is greater than %v", target)))
}

// IsEqualTo return success if the two values passed in are equal
func IsEqualTo(field string, value1 string, value2 string) ValidationResult {
	if value1 == value2 {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, "field values are not equal"))
}

// SatisfiesRegex return success if the field entered satifies the passed in regex
func SatisfiesRegex(field string, re string) ValidationResult {
	regex := regexp.MustCompile(re)
	if regex.MatchString(field) {
		return NewSuccessValidationResult()
	}

	return NewFailingValidationResult(NewValidationMessage(field, "field does not match regex"))
}
