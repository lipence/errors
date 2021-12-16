package errors

import (
	"bytes"
	"encoding/json"
	sysErr "errors"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"

	"go.uber.org/atomic"
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

type _packetASCType = struct {
	packageId uint32
	codeIdASC uint32
}

var _packetASCCount = atomic.Uint32{}
var _packetASCMap = sync.Map{} // used to store package asc id & code asc id
var _codeFileMap = sync.Map{}  // used to avoid redeclared code
var _mainPackagePath []byte    // used to trim errors.New caller package path
var _exprInternalCodeMask = regexp.MustCompile("^@([0-9a-f]{4})$")

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		_mainPackagePath = append([]byte(filepath.Clean(info.Main.Path)), os.PathSeparator)
	}
}

func New(code, message string) *underlying {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("failed to get caller package")
	}
	var packageASC *_packetASCType
	var packageCRC = crc32.ChecksumIEEE(bytes.TrimPrefix([]byte(filepath.Dir(file)), _mainPackagePath))
	{
		if _packetASCItf, ok := _packetASCMap.Load(packageCRC); ok {
			packageASC = _packetASCItf.(*_packetASCType)
		}
		if packageASC == nil {
			packageASC = &_packetASCType{_packetASCCount.Add(1), 0}
		}
		defer _packetASCMap.Store(packageCRC, packageASC)
	}
	{
		if code == "" {
			packageASC.codeIdASC += 1
			code = fmt.Sprintf("%08x%04x%04x", packageCRC, packageASC.packageId, packageASC.codeIdASC)
			log.Printf("empty error code at %s:%d , temporarily use code `%s` for message `%s`", file, line, code, message)
		} else if match := _exprInternalCodeMask.FindAllStringSubmatch(code, -1); len(match) > 0 {
			var codeLocalId = match[0][1]
			if codeLocalIdUint, err := strconv.ParseUint(codeLocalId, 16, 32); err != nil {
				panic(fmt.Errorf("invalid local code id: %v", err))
			} else if uint64(packageASC.codeIdASC) < codeLocalIdUint {
				packageASC.codeIdASC = uint32(codeLocalIdUint)
			}
			code = fmt.Sprintf("%08x%04x%04s", packageCRC, packageASC.packageId, codeLocalId)
		}
	}
	if anotherFile, ok := _codeFileMap.Load(code); ok {
		panic(fmt.Errorf("redeclared error code `%s`\n\tat: %s:%d\n\tat: %v", code, file, line, anotherFile))
	}
	_codeFileMap.Store(code, fmt.Sprintf("%s:%d", file, line))
	return &underlying{code: code, message: message}
}

func Because(underlying *underlying, cause error, fields ...Field) *node {
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
	case *underlying:
		return s
	case *node:
		return s
	case error:
		return s
	default:
		return NewSysErrf("%v", src)
	}
}

func Note(err error, fields ...Field) *node {
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

func Is(src error, target error) bool {
	if c, ok := src.(ComparableErr); ok {
		return c.Is(target)
	}
	return sysErr.Is(src, target)
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

type JsonErr struct {
	error
}

func (e JsonErr) MarshalJSON() ([]byte, error) {
	if e.error == nil {
		return []byte("null"), nil
	}
	return json.Marshal(e.Error())
}

func JSONMarshallable(err error) error {
	if _, ok := err.(json.Marshaler); ok {
		return err
	}
	return JsonErr{err}
}
