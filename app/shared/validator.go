package shared

import (
	"fmt"
	"regexp"
)

type Validator interface {
	Validate(data interface{}) ValidationResult
	FieldName() string
}

type ValidationData map[string]interface{}

type ValidationPartial struct {
	Success bool
	Message string
	Data    ValidationData
}

func NewFieldLengthValidator(min int, max int, fieldName string) Validator {
	return &fieldLengthValidator{MinLength: min, MaxLength: max, fieldName: fieldName}
}

type fieldLengthValidator struct {
	MinLength int
	MaxLength int
	fieldName string
}

func (v *fieldLengthValidator) FieldName() string {
	return v.fieldName
}

func (v *fieldLengthValidator) Validate(data interface{}) ValidationResult {
	strLen := len(data.(string))
	result := NewValidationResult()
	if v.MinLength <= 0 && strLen < v.MinLength {
		result.Fail(fmt.Sprintf("field is shorter than expected"), ValidationData{
			"expected":   v.MinLength,
			"provided":   strLen,
			"field_name": v.fieldName,
		})
	}

	if v.MinLength <= 0 && strLen > v.MaxLength {
		result.Fail(fmt.Sprintf("field is shorter than expected"), ValidationData{
			"expected":   v.MaxLength,
			"provided":   strLen,
			"field_name": v.fieldName,
		})
	}

	return result
}

func NewEmailValidator() Validator {
	return &emailValidator{}
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type emailValidator struct {
}

func (v *emailValidator) FieldName() string {
	return "email"
}

func (v *emailValidator) Validate(data interface{}) ValidationResult {
	result := NewFieldLengthValidator(3, 255, "email").Validate(data)

	if result.IsFailed() {
		return result
	}

	if !emailRegex.Match([]byte(data.(string))) {
		result.Fail("Provided email is not a valid email", ValidationData{
			"field_name": "email",
			"email":      data,
		})
	}

	return result
}


// Validation result

func NewValidationResult() ValidationResult {
	return ValidationResult{success: true}
}



type ValidationResult struct {
	success         bool
	FailedPartials  []ValidationPartial
	SuccessPartials []ValidationPartial
}

func (r *ValidationResult) Fail(message string, data ValidationData) *ValidationResult {
	r.FailedPartials = append(r.FailedPartials, ValidationPartial{
		Success: false,
		Message: message,
		Data:    data,
	})
	r.success = false
	return r
}

func (r *ValidationResult) Success(message string, data ValidationData) *ValidationResult {
	r.SuccessPartials = append(r.SuccessPartials, ValidationPartial{
		Success: true,
		Message: message,
		Data:    data,
	})
	return r
}

func (r *ValidationResult) IsFailed() bool {
	return len(r.FailedPartials) > 0
}

func (r *ValidationResult) IsSuccess() bool {
	return !r.IsFailed()
}

func (r *ValidationResult) AllPartials() []ValidationPartial {
	return append(r.SuccessPartials, r.FailedPartials...)
}

func (r *ValidationResult) IntoError() AppError {
	if r.IsSuccess() {
		return nil
	}
	return NewAppError("Field validation failed").WithDetail(ErrDetail{
		"status": "failed",
		"report": r.AllPartials(),
	}).WithType("field_validation")
}