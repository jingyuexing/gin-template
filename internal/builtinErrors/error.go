package builtinerrors

import "template/i18n"

type Exception struct {
	message      string
	path         string
	errorMessage map[string]string
	status       int
	code         int
}

func (e *Exception) Error() string {
	return e.message
}

func (e *Exception) Format(format ...map[string]any) string {
	e.message = i18n.I18N.T(e.path, format...)
	return e.message
}

func (e *Exception) ErrorMessage() map[string]string {
	return e.errorMessage
}

func (e *Exception) SetErrorMessage(field string, message string) {
	e.errorMessage[field] = message
}

func (e *Exception) SetMessage(message string) {
	e.message = message
}

func (e *Exception) Code() int {
	return e.code
}

func (e *Exception) Status() int {
	return e.status
}

func New(path string, status int, code int) *Exception {
	return &Exception{message: i18n.I18N.T(path), errorMessage: make(map[string]string), path: path, status: status, code: code}
}
