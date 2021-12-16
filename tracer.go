package errors

import (
	"fmt"
	"runtime"
)

type Tracer interface {
	Stack() []uintptr
}

type tracer struct {
	stack []uintptr
}

func (t *tracer) Stack() []uintptr {
	return t.stack
}

func (t *tracer) trace(skip int) {
	pcs := make([]uintptr, maxStackDepth)
	n := runtime.Callers(skip+2, pcs)
	t.stack = pcs[:n]
}

type traceInfoItem struct {
	Func string `json:"func"`
	Line string `json:"line"`
}

func (t *tracer) InfoStack(parent Tracer) []traceInfoItem {
	var currentStack = t.stack
	var sameFrames int
	// trim same stack frames
	if parent != nil {
		parentStack := parent.Stack()
	e1:
		for pi, pf := range parentStack {
			for ti, tf := range currentStack[:len(currentStack)-1-(len(parentStack)-1-pi)+1] {
				if pf == tf {
					sameFrames = len(currentStack) - (ti + 1) + 1
					break e1
				}
			}
		}
	}
	// build message
	var infoStack = make([]traceInfoItem, 0, len(currentStack)-sameFrames)
	frames := runtime.CallersFrames(currentStack)
	var more bool
	var frame runtime.Frame
	for i := len(currentStack) - 1 - sameFrames; i >= 0; i-- {
		frame, more = frames.Next()
		infoStack = append(infoStack, traceInfoItem{
			Func: fmt.Sprintf("[%d] %s", sameFrames+i, frame.Function),
			Line: fmt.Sprintf("%s:%d", frame.File, frame.Line),
		})
		if !more {
			break
		}
	}
	return infoStack
}
