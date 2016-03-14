package lorg

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func getReplacementsValues(replacements []replacement) []string {
	values := []string{}
	for _, replacement := range replacements {
		values = append(values, replacement.value)
	}

	return values
}

func getFilenameAndLine() string {
	callerPtr, file, line, _ := runtime.Caller(1)

	// caller will be like as following:
	// "github.com/kovetskiy/lorg.TestLogCallsFormatRender"
	caller := runtime.FuncForPC(callerPtr).Name()

	callerParts := strings.Split(caller, ".")
	funcname := callerParts[len(callerParts)-1]

	return fmt.Sprintf("[%s:%d] %s", filepath.Base(file), line, funcname)
}
