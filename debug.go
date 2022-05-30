//go:build debug
// +build debug

package errors

import (
	"bytes"
	"fmt"
	"hash/crc32"
	sysLog "log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"

	"go.uber.org/atomic"
)

type err_packetASCType = struct {
	packageId uint32
	codeIdASC uint32
}

var err_packetASCCount = atomic.Uint32{}
var err_packetASCMap = sync.Map{} // used to store package asc id & code asc id
var err_codeFileMap = sync.Map{}  // used to avoid redeclared code
var err_mainPackagePath []byte    // used to trim errors.New caller package path
var err_exprInternalCodeMask = regexp.MustCompile("^@([0-9a-f]{4})$")

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		err_mainPackagePath = append([]byte(filepath.Clean(info.Main.Path)), os.PathSeparator)
	}
}

func errMessageFilter(code, msg string) (newCode, newMessage string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("failed to get caller package")
	}
	var packageASC *err_packetASCType
	var packageCRC = crc32.ChecksumIEEE(bytes.TrimPrefix([]byte(filepath.Dir(file)), err_mainPackagePath))
	{
		if _packetASCItf, ok := err_packetASCMap.Load(packageCRC); ok {
			packageASC = _packetASCItf.(*err_packetASCType)
		}
		if packageASC == nil {
			packageASC = &err_packetASCType{err_packetASCCount.Add(1), 0}
		}
		defer err_packetASCMap.Store(packageCRC, packageASC)
	}
	{
		if code == "" {
			packageASC.codeIdASC += 1
			code = fmt.Sprintf("%08x%04x%04x", packageCRC, packageASC.packageId, packageASC.codeIdASC)
			sysLog.Printf("empty error code at %s:%d , temporarily use code `%s` for message `%s`", file, line, code, msg)
		} else if match := err_exprInternalCodeMask.FindAllStringSubmatch(code, -1); len(match) > 0 {
			var codeLocalId = match[0][1]
			if codeLocalIdUint, err := strconv.ParseUint(codeLocalId, 16, 32); err != nil {
				panic(fmt.Errorf("invalid local code id: %w", err))
			} else if uint64(packageASC.codeIdASC) < codeLocalIdUint {
				packageASC.codeIdASC = uint32(codeLocalIdUint)
			}
			code = fmt.Sprintf("%08x%04x%04s", packageCRC, packageASC.packageId, codeLocalId)
		}
	}
	if anotherFile, ok := err_codeFileMap.Load(code); ok {
		panic(fmt.Errorf("redeclared error code `%s`\n\tat: %s:%d\n\tat: %v", code, file, line, anotherFile))
	}
	err_codeFileMap.Store(code, fmt.Sprintf("%s:%d", file, line))
	return code, msg
}

func init() {
	OnCreateMsg(errMessageFilter)
}
