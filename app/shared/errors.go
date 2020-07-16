package shared

import (
	log "github.com/sirupsen/logrus"
)

const (
	ErrStatusInvalidRequest string = "invalid_request"
	ErrStatusNotFound       string = "not_found"
)

type ErrDetail map[string]interface{}

func NewGouthyError(msg string) GouthyError {
	mixin := gouthyErrorImpl{
		message: msg,
		detail:  make(map[string]interface{}),
		errType: "general",
	}
	return &mixin
}

type GouthyError interface {
	error
	WithDetail(detail ErrDetail) GouthyError
	AddDetailField(name string, value interface{}) GouthyError
	WithType(errType string) GouthyError
	Type() string
	Detail() ErrDetail
	Message() string
	LogAppend(entry *log.Entry) *log.Entry
}

type gouthyErrorImpl struct {
	message string
	detail  ErrDetail
	errType string
}

func (err *gouthyErrorImpl) Type() string {
	return err.errType
}

func (err *gouthyErrorImpl) WithType(errType string) GouthyError {
	err.errType = errType
	return err
}

func (err *gouthyErrorImpl) AddDetailField(name string, value interface{}) GouthyError {
	err.detail[name] = value
	return err
}

func (err *gouthyErrorImpl) WithDetail(partial ErrDetail) GouthyError {
	for key, val := range partial {
		err.detail[key] = val
	}
	return err
}

func (err *gouthyErrorImpl) Detail() ErrDetail {
	return err.detail
}

func (err *gouthyErrorImpl) Message() string {
	return err.message
}

func (err *gouthyErrorImpl) Error() string {
	return err.message
}

func (err *gouthyErrorImpl) LogAppend(entry *log.Entry) *log.Entry {
	return entry.WithError(err).WithFields(log.Fields(err.detail))
}

func LogError(entry *log.Entry, err error) *log.Entry {
	switch v := err.(type) {
	case GouthyError:
		return v.LogAppend(entry)
	default:
		return entry.WithError(err)
	}
}
