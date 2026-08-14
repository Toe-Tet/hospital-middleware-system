package errors

const (
	CodeBadRequest          = "BAD_REQUEST"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeForbidden           = "FORBIDDEN"
	CodeNotFound            = "NOT_FOUND"
	CodeConflict            = "CONFLICT"
	CodeValidationError     = "VALIDATION_ERROR"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
)

const (
	MsgBadRequest          = "Bad request"
	MsgUnauthorized        = "Unauthorized"
	MsgForbidden           = "Forbidden"
	MsgNotFound            = "Resource not found"
	MsgConflict            = "Resource already exists"
	MsgValidationError     = "Validation failed"
	MsgInternalServerError = "Internal server error"
)
