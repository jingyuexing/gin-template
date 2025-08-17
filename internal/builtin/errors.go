package builtin

import (
	"net/http"
	builtinerrors "template/internal/builtinErrors"
)

// System level error codes (1000-1999)
const (
	ParameterIsInvalidCode = 1000 + iota
	ServerErrorCode
	RateLimitExceededCode
	NotFoundCode
	MethodNotAllowedCode
	RequestTimeoutCode
	ConflictCode
	TooManyRequestsCode
	FileNotFoundCode
	FileTooLargeCode
	LimitRateCode
)

// Authentication & Authorization error codes (2000-2999)
const (
	UnauthorizedCode = 2000 + iota
	ForbiddenCode
	UserNotFoundCode
	UserPasswordIncorrectCode
	UserNotVerifiedCode
	UserNameExistsCode
	TokenExpiredCode
	TokenInvalidCode
	TokenRequiredCode    // 新增: Token必需错误码
	InvalidPasswordCode  // 新增: 密码无效错误码
	TokenTypeInvalidCode // 新增: Token类型无效错误码
)

// Database operation error codes (3000-3999)
const (
	DBQueryFailedCode = 3000 + iota
	DBRecordNotFoundCode
	DBInsertFailedCode
	DBDeleteFailedCode
	DBUpdateFailedCode
)

// System level errors
var (
	ErrInternalServer   = builtinerrors.New("Errors.System.ServerError", http.StatusInternalServerError, ServerErrorCode)
	ErrRateLimit        = builtinerrors.New("Errors.System.RateLimit", http.StatusTooManyRequests, RateLimitExceededCode)
	ErrNotFound         = builtinerrors.New("Errors.System.NotFound", http.StatusNotFound, NotFoundCode)
	ErrMethodNotAllowed = builtinerrors.New("Errors.System.MethodNotAllowed", http.StatusMethodNotAllowed, MethodNotAllowedCode)
	ErrTimeout          = builtinerrors.New("Errors.System.Timeout", http.StatusRequestTimeout, RequestTimeoutCode)
	ErrConflict         = builtinerrors.New("Errors.System.Conflict", http.StatusConflict, ConflictCode)
	ErrTooManyRequests  = builtinerrors.New("Errors.System.TooManyRequests", http.StatusTooManyRequests, TooManyRequestsCode)
	ErrFileNotFound     = builtinerrors.New("Errors.System.FileNotFound", http.StatusNotFound, FileNotFoundCode)
	ErrFileTooLarge     = builtinerrors.New("Errors.System.FileTooLarge", http.StatusRequestEntityTooLarge, FileTooLargeCode)
	ErrInvalidParams    = builtinerrors.New("Errors.System.InvalidParams", http.StatusBadRequest, ParameterIsInvalidCode)
	ErrRequestLimitRate = builtinerrors.New("Errors.System.LimitRate", http.StatusTooManyRequests, LimitRateCode)
	// Database operation errors
	ErrDBQueryFailed    = builtinerrors.New("Errors.DB.QueryFailed", http.StatusInternalServerError, DBQueryFailedCode)
	ErrDBRecordNotFound = builtinerrors.New("Errors.DB.RecordNotFound", http.StatusNotFound, DBRecordNotFoundCode)
	ErrDBInsertFailed   = builtinerrors.New("Errors.DB.InsertFailed", http.StatusInternalServerError, DBInsertFailedCode)
	ErrDBDeleteFailed   = builtinerrors.New("Errors.DB.DeleteFailed", http.StatusInternalServerError, DBDeleteFailedCode)
	ErrDBUpdateFailed   = builtinerrors.New("Errors.DB.UpdateFailed", http.StatusInternalServerError, DBUpdateFailedCode)
)

// Authentication & Authorization errors
var (
	ErrUnauthorized          = builtinerrors.New("Errors.Auth.Unauthorized", http.StatusUnauthorized, UnauthorizedCode)
	ErrForbidden             = builtinerrors.New("Errors.Auth.Forbidden", http.StatusForbidden, ForbiddenCode)
	ErrUserNotFound          = builtinerrors.New("Errors.Auth.UserNotFound", http.StatusNotFound, UserNotFoundCode)
	ErrUserPasswordIncorrect = builtinerrors.New("Errors.Auth.PasswordIncorrect", http.StatusUnauthorized, UserPasswordIncorrectCode)
	ErrUserNotVerified       = builtinerrors.New("Errors.Auth.UserNotVerified", http.StatusForbidden, UserNotVerifiedCode)
	ErrUserNameExists        = builtinerrors.New("Errors.Auth.UserNameExists", http.StatusConflict, UserNameExistsCode)
	ErrTokenExpired          = builtinerrors.New("Errors.Auth.TokenExpired", http.StatusUnauthorized, TokenExpiredCode)
	ErrTokenInvalid          = builtinerrors.New("Errors.Auth.TokenInvalid", http.StatusUnauthorized, TokenInvalidCode)
	// 新增以下错误定义
	ErrTokenRequired    = builtinerrors.New("Errors.Auth.TokenRequired", http.StatusUnauthorized, TokenRequiredCode)
	ErrInvalidPassword  = builtinerrors.New("Errors.Auth.InvalidPassword", http.StatusUnauthorized, InvalidPasswordCode)
	ErrTokenTypeInvalid = builtinerrors.New("Errors.Auth.TokenTypeInvalid", http.StatusUnauthorized, TokenTypeInvalidCode)
)
