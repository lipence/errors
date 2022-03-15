package errors

import (
	"fmt"
	"time"
	_ "unsafe"

	_ "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	_ "go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

//go:linkname Any go.uber.org/zap.Any
func Any(key string, value interface{}) Field

//go:linkname Array go.uber.org/zap.Array
func Array(key string, val zapcore.ArrayMarshaler) Field

//go:linkname Binary go.uber.org/zap.Binary
func Binary(key string, val []byte) Field

//go:linkname Bool go.uber.org/zap.Bool
func Bool(key string, val bool) Field

//go:linkname Boolp go.uber.org/zap.Boolp
func Boolp(key string, val *bool) Field

//go:linkname Bools go.uber.org/zap.Bools
func Bools(key string, bs []bool) Field

//go:linkname ByteString go.uber.org/zap.ByteString
func ByteString(key string, val []byte) Field

//go:linkname ByteStrings go.uber.org/zap.ByteStrings
func ByteStrings(key string, bss [][]byte) Field

//go:linkname Complex128 go.uber.org/zap.Complex128
func Complex128(key string, val complex128) Field

//go:linkname Complex128p go.uber.org/zap.Complex128p
func Complex128p(key string, val *complex128) Field

//go:linkname Complex128s go.uber.org/zap.Complex128s
func Complex128s(key string, nums []complex128) Field

//go:linkname Complex64 go.uber.org/zap.Complex64
func Complex64(key string, val complex64) Field

//go:linkname Complex64p go.uber.org/zap.Complex64p
func Complex64p(key string, val *complex64) Field

//go:linkname Complex64s go.uber.org/zap.Complex64s
func Complex64s(key string, nums []complex64) Field

//go:linkname Duration go.uber.org/zap.Duration
func Duration(key string, val time.Duration) Field

//go:linkname Durationp go.uber.org/zap.Durationp
func Durationp(key string, val *time.Duration) Field

//go:linkname Durations go.uber.org/zap.Durations
func Durations(key string, ds []time.Duration) Field

//go:linkname Error go.uber.org/zap.Error
func Error(err error) Field

//go:linkname Errors go.uber.org/zap.Errors
func Errors(key string, errs []error) Field

//go:linkname Float32 go.uber.org/zap.Float32
func Float32(key string, val float32) Field

//go:linkname Float32p go.uber.org/zap.Float32p
func Float32p(key string, val *float32) Field

//go:linkname Float32s go.uber.org/zap.Float32s
func Float32s(key string, nums []float32) Field

//go:linkname Float64 go.uber.org/zap.Float64
func Float64(key string, val float64) Field

//go:linkname Float64p go.uber.org/zap.Float64p
func Float64p(key string, val *float64) Field

//go:linkname Float64s go.uber.org/zap.Float64s
func Float64s(key string, nums []float64) Field

//go:linkname Inline go.uber.org/zap.Inline
func Inline(val zapcore.ObjectMarshaler) Field

//go:linkname Int go.uber.org/zap.Int
func Int(key string, val int) Field

//go:linkname Int16 go.uber.org/zap.Int16
func Int16(key string, val int16) Field

//go:linkname Int16p go.uber.org/zap.Int16p
func Int16p(key string, val *int16) Field

//go:linkname Int16s go.uber.org/zap.Int16s
func Int16s(key string, nums []int16) Field

//go:linkname Int32 go.uber.org/zap.Int32
func Int32(key string, val int32) Field

//go:linkname Int32p go.uber.org/zap.Int32p
func Int32p(key string, val *int32) Field

//go:linkname Int32s go.uber.org/zap.Int32s
func Int32s(key string, nums []int32) Field

//go:linkname Int64 go.uber.org/zap.Int64
func Int64(key string, val int64) Field

//go:linkname Int64p go.uber.org/zap.Int64p
func Int64p(key string, val *int64) Field

//go:linkname Int64s go.uber.org/zap.Int64s
func Int64s(key string, nums []int64) Field

//go:linkname Int8 go.uber.org/zap.Int8
func Int8(key string, val int8) Field

//go:linkname Int8p go.uber.org/zap.Int8p
func Int8p(key string, val *int8) Field

//go:linkname Int8s go.uber.org/zap.Int8s
func Int8s(key string, nums []int8) Field

//go:linkname Intp go.uber.org/zap.Intp
func Intp(key string, val *int) Field

//go:linkname Ints go.uber.org/zap.Ints
func Ints(key string, nums []int) Field

//go:linkname NamedError go.uber.org/zap.NamedError
func NamedError(key string, err error) Field

//go:linkname Namespace go.uber.org/zap.Namespace
func Namespace(key string) Field

//go:linkname Object go.uber.org/zap.Object
func Object(key string, val zapcore.ObjectMarshaler) Field

//go:linkname Reflect go.uber.org/zap.Reflect
func Reflect(key string, val interface{}) Field

//go:linkname Skip go.uber.org/zap.Skip
func Skip() Field

//go:linkname Stack go.uber.org/zap.Stack
func Stack(key string) Field

//go:linkname StackSkip go.uber.org/zap.StackSkip
func StackSkip(key string, skip int) Field

//go:linkname String go.uber.org/zap.String
func String(key string, val string) Field

//go:linkname Stringer go.uber.org/zap.Stringer
func Stringer(key string, val fmt.Stringer) Field

//go:linkname Stringp go.uber.org/zap.Stringp
func Stringp(key string, val *string) Field

//go:linkname Strings go.uber.org/zap.Strings
func Strings(key string, ss []string) Field

//go:linkname Time go.uber.org/zap.Time
func Time(key string, val time.Time) Field

//go:linkname Timep go.uber.org/zap.Timep
func Timep(key string, val *time.Time) Field

//go:linkname Times go.uber.org/zap.Times
func Times(key string, ts []time.Time) Field

//go:linkname Uint go.uber.org/zap.Uint
func Uint(key string, val uint) Field

//go:linkname Uint16 go.uber.org/zap.Uint16
func Uint16(key string, val uint16) Field

//go:linkname Uint16p go.uber.org/zap.Uint16p
func Uint16p(key string, val *uint16) Field

//go:linkname Uint16s go.uber.org/zap.Uint16s
func Uint16s(key string, nums []uint16) Field

//go:linkname Uint32 go.uber.org/zap.Uint32
func Uint32(key string, val uint32) Field

//go:linkname Uint32p go.uber.org/zap.Uint32p
func Uint32p(key string, val *uint32) Field

//go:linkname Uint32s go.uber.org/zap.Uint32s
func Uint32s(key string, nums []uint32) Field

//go:linkname Uint64 go.uber.org/zap.Uint64
func Uint64(key string, val uint64) Field

//go:linkname Uint64p go.uber.org/zap.Uint64p
func Uint64p(key string, val *uint64) Field

//go:linkname Uint64s go.uber.org/zap.Uint64s
func Uint64s(key string, nums []uint64) Field

//go:linkname Uint8 go.uber.org/zap.Uint8
func Uint8(key string, val uint8) Field

//go:linkname Uint8p go.uber.org/zap.Uint8p
func Uint8p(key string, val *uint8) Field

//go:linkname Uint8s go.uber.org/zap.Uint8s
func Uint8s(key string, nums []uint8) Field

//go:linkname Uintp go.uber.org/zap.Uintp
func Uintp(key string, val *uint) Field

//go:linkname Uintptr go.uber.org/zap.Uintptr
func Uintptr(key string, val uintptr) Field

//go:linkname Uintptrp go.uber.org/zap.Uintptrp
func Uintptrp(key string, val *uintptr) Field

//go:linkname Uintptrs go.uber.org/zap.Uintptrs
func Uintptrs(key string, us []uintptr) Field

//go:linkname Uints go.uber.org/zap.Uints
func Uints(key string, nums []uint) Field
