package shared

import (
	"fmt"
)

func CheckFieldMinLength(field string, fieldName string, minLen int) AppError {
	fieldLen := len(field)
	if fieldLen < minLen {
		return NewErrInvalidField(fieldName).WithDetail(ErrDetail{
			"reason":     fmt.Sprintf("%s is too short", fieldName),
			"min_length": minLen,
			"length":     fieldLen,
		})
	}
	return nil
}

func CheckFieldMaxLength(field string, fieldName string, max int) AppError {
	fieldLen := len(field)
	if fieldLen > max {
		return NewErrInvalidField(fieldName).WithDetail(ErrDetail{
			"reason":     fmt.Sprintf("%s is too long", fieldName),
			"max_length": max,
			"length":     fieldLen,
		})
	}
	return nil
}

