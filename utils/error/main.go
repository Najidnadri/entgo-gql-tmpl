package apperror

import (
	"fmt"
	"strings"
)

type AppError struct {
	Message string   `json:"message"`
	Stack   []string `json:"stack"`
	Err     error    `json:"error"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewError(err error, message string, replaceMsg bool) error {
	appErr, ok := err.(AppError)
	if !ok {
		return AppError{
			Message: message,
			Err:     err,
			Stack:   []string{fmt.Sprintf("%s: %v", message, err)},
		}
	} else {
		if replaceMsg {
			appErr.Message = message
		}
		appErr.Stack = append([]string{message}, appErr.Stack...)
		return appErr
	}
}

func GetStack(err error) string {
	appErr, ok := err.(AppError)
	if !ok {
		return err.Error()
	}
	return fmt.Sprintf(`[err:%v] [msg:%s] stack:%s`, appErr.Err, appErr.Message, strings.Join(appErr.Stack, " . "))
}
