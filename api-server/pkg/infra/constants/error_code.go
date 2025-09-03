package constants

type ErrorCode string

const (
	// CodeRequired General errors
	CodeRequired      = ErrorCode("PARAM.REQUIRED")
	CodeInvalidParam  = ErrorCode("PARAM.INVALID.PARAM")
	CodeInvalidScheme = ErrorCode("PARAM.INVALID.SCHEME")
	CodeNotListParam  = ErrorCode("PARAM.NOT_LIST")
	CodeInvalidJson   = ErrorCode("PARAM.INVALID.JSON")
	CodeInvalidData   = ErrorCode("PARAM.INVALID.DATA")

	CodeInternalError            = ErrorCode("ERROR.INTERNAL")
	CodeBackingImageCreatedError = ErrorCode("ERROR.BACKINGIMAGE.CREATED.FAILED")

	CodeValidationFailed = ErrorCode("PARAM.VALIDATION.FAILED")
)
