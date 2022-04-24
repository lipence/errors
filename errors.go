package errors

import (
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
	if target == nil {
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
