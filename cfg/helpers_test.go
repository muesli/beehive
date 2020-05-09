package cfg

import (
	"runtime"
	"strings"
)

// Replace backward slashes in Windows paths with /, to make them suitable
// for Go URL parsing.
func fixWindowsPath(path string) string {
	if runtime.GOOS == "windows" {
		return strings.Replace(path, `\`, "/", -1)
	}

	return path
}
