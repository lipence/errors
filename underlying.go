package errors

import "encoding/json"

type Message interface {
	Code() string
	Message() string
}

type underlying struct {
	code    string
	message string
}

func (e *underlying) Code() string {
	return e.code
}

func (e *underlying) Message() string {
	return e.message
}

func (e *underlying) Is(target error) bool {
	if u, ok := target.(*underlying); ok {
		if e.code != "" {
			return e.code == u.code
		}
		if e.message != "" {
			return e.message == u.message
		}
		return false
	}
	return Is(target, e)
}

func (e *underlying) Error() string {
	if e == nil {
		return ""
	}
	return e.code + ": " + e.message
}

func (e *underlying) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error())
}

func NewUnderlying(code, message string) *underlying {
	return &underlying{code: code, message: message}
}
