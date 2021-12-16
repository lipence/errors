package errors

import (
	"bytes"
	"encoding/json"
	"sync"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type node struct {
	tracer
	data       []Field
	dataRWM    sync.RWMutex
	underlying *underlying
	cause      error
}

func (e *node) Is(target error) bool {
	if targetNode, ok := target.(*node); ok {
		return Is(e.Underlying(), targetNode.Underlying())
	}
	return Is(e.Underlying(), target)
}

func (e *node) CausedBy(target error, deepFirst bool) (bool, error) {
	if e.underlying != nil {
		if Is(e.underlying, target) {
			return true, e
		}
	}
	if e.cause != nil {
		if !deepFirst && Is(e.cause, target) {
			return true, e
		}
		if n, ok := e.cause.(*node); ok {
			return n.CausedBy(target, deepFirst)
		}
		if deepFirst && Is(e.cause, target) {
			return true, e
		}
	}
	return false, nil
}

func (e *node) Error() string {
	return e.message()
}

func (e *node) Code() (code string) {
	if e.cause != nil {
		if cm, ok := e.cause.(Message); ok {
			if code = cm.Code(); code != "" {
				return code
			}
		}
	}
	if e.underlying != nil {
		if code = e.underlying.Code(); code != "" {
			return code
		}
	}
	return ""
}

func (e *node) Message() (message string) {
	if e.cause != nil {
		if cm, ok := e.cause.(Message); ok {
			if message = cm.Message(); message != "" {
				return message
			}
		}
	}
	if e.underlying != nil {
		if message = e.underlying.Message(); message != "" {
			return message
		}
	}
	return ""
}

func (e *node) WithData(fields ...Field) *node {
	if e == nil {
		return nil
	}
	e.dataRWM.Lock()
	e.data = append(e.data, fields...)
	e.dataRWM.Unlock()
	return e
}

func (e *node) DataMap() map[string]interface{} {
	var me = zapcore.NewMapObjectEncoder()
	e.dataRWM.RLock()
	for i := 0; i < len(e.data); i++ {
		e.data[i].AddTo(me)
	}
	e.dataRWM.RUnlock()
	return me.Fields
}

func (e *node) Data(key string, r bool) (val interface{}, found bool) {
	if r {
		if causeNode, ok := e.cause.(*node); ok {
			if val, found = causeNode.Data(key, r); found {
				return val, true
			}
		}
	}
	for i := 0; i < len(e.data); i++ {
		if e.data[i].Key != key {
			continue
		}
		var me = zapcore.NewMapObjectEncoder()
		e.data[i].AddTo(me)
		return me.Fields[key], true
	}
	return nil, false
}

func (e *node) HasData(key string, r bool) (found bool) {
	if r {
		if causeNode, ok := e.cause.(*node); ok {
			if _, found = causeNode.Data(key, r); found {
				return true
			}
		}
	}
	for i := 0; i < len(e.data); i++ {
		if e.data[i].Key == key {
			return true
		}
	}
	return false
}

func (e *node) Underlying() error {
	if e.underlying != nil {
		return e.underlying
	}
	if e.cause != nil {
		if _, ok := e.cause.(*node); !ok {
			return e.cause
		}
	}
	return nil
}

func (e *node) Unwrap() error {
	return e.cause
}

func (e *node) Cause() error {
	return e.cause
}

func (e *node) WithCause(err error) *node {
	e.cause = err
	return e
}

type nodeData []Field

func (d nodeData) MarshalJSON() (dst []byte, err error) {
	var b *buffer.Buffer
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{})
	if b, err = encoder.EncodeEntry(zapcore.Entry{}, d); err == nil && b != nil {
		dst = b.Bytes()
	}
	return dst, err
}

type nodeInfoItem struct {
	Underlying error           `json:"underlying"`
	Data       nodeData        `json:"data,omitempty"`
	StackTrace []traceInfoItem `json:"stackTrace,omitempty"`
}

func (e *node) InfoStack(parent *node) []nodeInfoItem {
	var stack []nodeInfoItem
	var nodeItem nodeInfoItem
	if e.cause != nil {
		if causeNode, ok := e.cause.(*node); ok {
			stack = causeNode.InfoStack(e)
		} else if e.underlying != nil {
			stack = append(stack, nodeInfoItem{
				Underlying: JSONMarshallable(e.cause),
			})
		} else {
			nodeItem = nodeInfoItem{Underlying: JSONMarshallable(e.cause)}
		}
	}
	if e.underlying != nil {
		nodeItem = nodeInfoItem{Underlying: JSONMarshallable(e.underlying)}
	}
	e.dataRWM.RLock()
	if len(e.data) > 0 {
		nodeItem.Data = e.data
	}
	e.dataRWM.RUnlock()
	if parent != nil {
		nodeItem.StackTrace = e.tracer.InfoStack(&parent.tracer)
	} else {
		nodeItem.StackTrace = e.tracer.InfoStack(nil)
	}
	stack = append(stack, nodeItem)
	return stack
}

func (e *node) message() string {
	var b = bytes.Buffer{}
	for _, infoItem := range e.InfoStack(nil) {
		if infoItem.Underlying != nil {
			b.WriteRune('\n')
			b.WriteString(infoItem.Underlying.Error())
		}
		if len(infoItem.Data) > 0 {
			b.Write([]byte{':', '\x20'})
			if data, err := json.Marshal(infoItem.Data); err == nil {
				b.Write(data)
			} else {
				b.WriteString(err.Error())
			}
		}
		for _, traceItem := range infoItem.StackTrace {
			b.WriteString("\n\x20\x20")
			b.WriteString(traceItem.Func)
			b.WriteString("\n\x20\x20\x20\x20")
			b.WriteString(traceItem.Line)
		}
	}
	return b.String()
}

func (e *node) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.InfoStack(nil))
}
