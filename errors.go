package errors

import (
	"encoding/json"
	sysErr "errors"
	"fmt"
	"sync"
)

const maxStackDepth = 64

type ComparableErr interface {
	Is(error) bool
}

func NewSysErr(message string) error {
	return sysErr.New(message)
}

func NewSysErrf(message string, args ...interface{}) error {
	return fmt.Errorf(message, args...)
}

type MsgFilter func(code, msg string) (newCode, newMessage string)

var msgFilters []MsgFilter

func OnCreateMsg(filter MsgFilter) {
	if filter != nil {
		msgFilters = append(msgFilters, filter)
	}
}

func New(code, msg string) *underlying {
	for i := 0; i < len(msgFilters); i++ {
		if filter := msgFilters[i]; filter != nil {
			code, msg = filter(code, msg)
		}
	}
	return &underlying{code: code, message: msg}
}

func Because(underlying *underlying, cause error, fields ...Field) *node {
	if cause == nil {
		return nil
	}
	n := &node{
		dataRWM:    sync.RWMutex{},
		underlying: underlying,
		cause:      cause,
	}
	n.trace(1)
	return n.WithData(fields...)
}

func From(src interface{}) error {
	switch s := src.(type) {
	case nil:
		return nil
	case *node:
		return s
	case *underlying:
		return s
	case error:
		return s
	default:
		return NewSysErrf("%v", src)
	}
}

func Note(err error, fields ...Field) *node {
	if err == nil {
		return nil
	}
	if n, ok := err.(*node); ok {
		return n.WithData(fields...)
	}
	n := &node{
		dataRWM: sync.RWMutex{},
		cause:   err,
	}
	n.trace(1)
	return n.WithData(fields...)
}

func Is(err error, target error) bool {
	if err == nil || target == nil {
		return err == target
	}
	if c, ok := err.(ComparableErr); ok {
		return c.Is(target)
	}
	return sysErr.Is(err, target)
}

func Data(src error, key string, r bool) (interface{}, bool) {
	if srcNode, ok := src.(*node); ok {
		return srcNode.Data(key, r)
	}
	return nil, false
}

func HasData(src error, key string, r bool) bool {
	if srcNode, ok := src.(*node); ok {
		return srcNode.HasData(key, r)
	}
	return false
}

func CausedBy(src error, target error, deepFirst bool) bool {
	return CausedByNode(src, target, deepFirst, nil)
}

func CausedByNode(src error, target error, deepFirst bool, causeReceiver *error) (isCausedBy bool) {
	var cause error
	if n, ok := src.(*node); ok {
		isCausedBy, cause = n.CausedBy(target, deepFirst)
		if causeReceiver != nil {
			*causeReceiver = cause
		}
		return isCausedBy
	}
	if Is(src, target) {
		if causeReceiver != nil {
			*causeReceiver = src
		}
		return true
	}
	return false
}

type BatchErrors []error

func (e BatchErrors) Error() string {
	var buf = getBytesBuffer()
	defer returnBytesBuffer(buf)
	buf.WriteString("Multiple error occurred:")
	for i, err := range e {
		if err != nil {
			buf.WriteString(fmt.Sprintf("\n[%d] %v", i, err))
		}
	}
	return buf.String()
}

func (e BatchErrors) MarshalJSON() ([]byte, error) {
	var errArr = e
	for i, err := range errArr {
		errArr[i] = AsJsonMarshaller(err)
	}
	return json.Marshal(errArr)
}

func Batch(errs []error) error {
	var _errs = make([]error, 0, len(errs))
	for _, _err := range errs {
		if _err != nil {
			_errs = append(_errs, _err)
		}
	}
	switch len(_errs) {
	case 0:
		return nil
	case 1:
		return _errs[0]
	default:
		return BatchErrors(_errs)
	}
}

func Unbatch(errs error) ([]error, bool) {
	if errs == nil {
		return nil, true
	}
	errArr, ok := errs.(BatchErrors)
	if !ok {
		return nil, false
	}
	for i, err := range errArr {
		if je, isJsonErr := err.(jsonErr); isJsonErr {
			errArr[i] = je.error
		}
	}
	return errArr, true
}
