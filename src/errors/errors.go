package errors

import "fmt"

type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewBadRequest(msg string, details ...interface{}) *AppError {
	if msg == "" {
		msg = MsgBadRequest
	}
	err := &AppError{
		Code:    CodeBadRequest,
		Message: msg,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

func NewUnauthorized(msg string, details ...interface{}) *AppError {
	if msg == "" {
		msg = MsgUnauthorized
	}
	err := &AppError{
		Code:    CodeUnauthorized,
		Message: msg,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

func NewForbidden(msg string, details ...interface{}) *AppError {
	if msg == "" {
		msg = MsgForbidden
	}
	err := &AppError{
		Code:    CodeForbidden,
		Message: msg,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

func NewNotFound(msg string, details ...interface{}) *AppError {
	if msg == "" {
		msg = MsgNotFound
	}
	err := &AppError{
		Code:    CodeNotFound,
		Message: msg,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

func NewConflict(msg string, details ...interface{}) *AppError {
	if msg == "" {
		msg = MsgConflict
	}
	err := &AppError{
		Code:    CodeConflict,
		Message: msg,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

func NewValidationError(msg string, details interface{}) *AppError {
	if msg == "" {
		msg = MsgValidationError
	}
	return &AppError{
		Code:    CodeValidationError,
		Message: msg,
		Details: details,
	}
}

func NewInternal(err error, msg ...string) *AppError {
	message := MsgInternalServerError
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	return &AppError{
		Code:    CodeInternalServerError,
		Message: message,
		Err:     err,
	}
}

func StatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		switch appErr.Code {
		case CodeBadRequest, CodeValidationError:
			return 400
		case CodeUnauthorized:
			return 401
		case CodeForbidden:
			return 403
		case CodeNotFound:
			return 404
		case CodeConflict:
			return 409
		case CodeInternalServerError:
			return 500
		default:
			return 500
		}
	}
	return 500
}
