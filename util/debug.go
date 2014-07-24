package util

import (
	"reflect"
	"runtime"
	"strings"
)

//------------------------------------------------------------
// Debugging utils
//------------------------------------------------------------

// Discovers function info: name and package.
func GetFuncInfo(obj interface{}) []string {
	fullname := runtime.FuncForPC(reflect.ValueOf(obj).Pointer()).Name()
	idx := strings.LastIndex(fullname, ".")
	if idx == -1 {
		return []string{"-", "-"}
	}

	return []string{
		fullname[idx+1:],
		fullname[0:idx],
	}
}

// Returns full stack trace as string.
func StackFull() string {
	return stack(0)
}

// Returns stack trace as string, skipping few first lines.
func Stack() string {
	return stack(5)
}

// Stack implementation.
func stack(skip int) string {
	buf := make([]byte, 1024)
	size := 0
	for {
		size = runtime.Stack(buf, false)
		if size < len(buf) {
			break
		}
		buf = make([]byte, len(buf)*2)
	}

	// Skip first 5 newlines
	var count = 0
	var idx = 0
	for idx < size {
		if count == skip {
			break
		}
		if buf[idx] == '\n' {
			count++
		}
		idx++
	}
	return string(buf[idx:size])
}
