//go:build windows
// +build windows

package autocomplete

import (
	"strings"

	"github.com/lmorg/murex/utils/consts"
)

// SplitPath takes a %PATH% (or similar) string and splits it into an array of paths
func SplitPath(envPath string) []string {
	split := strings.Split(envPath, ";")
	return split
}

// listExes called listExesWindows which exists in execs.go because it needs to
// be called when murex runs inside WSL
func listExes(path string, exes map[string]bool) {
	listExesWindows(path, exes)
}

func matchExes(s string, exes map[string]bool, includeColon bool) (items []string) {
	colon := ""
	if includeColon {
		colon = ":"
	}

	for name := range exes {
		lc := strings.ToLower(s)
		if strings.HasPrefix(strings.ToLower(name), lc) {
			switch {
			case isSpecialBuiltin(name):
				items = append(items, name[len(s):]+colon)
			case consts.NamedPipeProcName == name:
				// do nothing
			default:
				items = append(items, name[len(s):])
			}
		}
	}

	return
}
