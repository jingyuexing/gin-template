package i18n

type SystemErrors struct {
	ServerError      string `json:"server_error" yaml:"server_error"`
	RateLimit        string `json:"rate_limit" yaml:"rate_limit"`
	NotFound         string `json:"not_found" yaml:"not_found"`
	MethodNotAllowed string `json:"method_not_allowed" yaml:"method_not_allowed"`
	Timeout          string `json:"timeout" yaml:"timeout"`
	Conflict         string `json:"conflict" yaml:"conflict"`
	TooManyRequests  string `json:"too_many_requests" yaml:"too_many_requests"`
	FileNotFound     string `json:"file_not_found" yaml:"file_not_found"`
	FileTooLarge     string `json:"file_too_large" yaml:"file_too_large"`
}

type AuthErrors struct {
	Unauthorized      string `json:"unauthorized" yaml:"unauthorized"`
	Forbidden         string `json:"forbidden" yaml:"forbidden"`
	UserNotFound      string `json:"user_not_found" yaml:"user_not_found"`
	PasswordIncorrect string `json:"password_incorrect" yaml:"password_incorrect"`
	UserNotVerified   string `json:"user_not_verified" yaml:"user_not_verified"`
	UserNameExists    string `json:"username_exists" yaml:"username_exists"`
	TokenExpired      string `json:"token_expired" yaml:"token_expired"`
	TokenInvalid      string `json:"token_invalid" yaml:"token_invalid"`
}

// 新增数据库错误定义
type DBErrors struct {
	QueryFailed    string `json:"query_failed" yaml:"query_failed"`
	RecordNotFound string `json:"record_not_found" yaml:"record_not_found"`
	InsertFailed   string `json:"insert_failed" yaml:"insert_failed"`
	DeleteFailed   string `json:"delete_failed" yaml:"delete_failed"`
	UpdateFailed   string `json:"update_failed" yaml:"update_failed"`
}

type Errors struct {
	System SystemErrors `json:"system" yaml:"system"`
	Auth   AuthErrors   `json:"auth" yaml:"auth"`
	DB     DBErrors     `json:"db" yaml:"db"`
}

type Logger struct {
}

type Locale struct {
	Errors Errors `json:"errors" yaml:"errors"`
	Logger Logger `json:"logger" yaml:"logger"`
}
