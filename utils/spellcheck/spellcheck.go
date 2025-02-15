package spellcheck

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/lmorg/murex/lang"
	"github.com/lmorg/murex/lang/types"
	"github.com/lmorg/murex/shell/autocomplete"
	"github.com/lmorg/murex/utils"
	"github.com/lmorg/murex/utils/ansi"
	"github.com/lmorg/murex/utils/spellcheck/userdictionary"
)

var rxRemoveEsc = regexp.MustCompile(`\\[a-zA-Z]`)

// String spellchecks a line of type string and returns an underlined (ANSI escaped) string
func String(line string) (string, error) {
	enabled, err := lang.ShellProcess.Config.Get("shell", "spellcheck-enabled", types.Boolean)
	if err != nil || !enabled.(bool) {
		return line, err
	}

	block, err := lang.ShellProcess.Config.Get("shell", "spellcheck-func", types.CodeBlock)
	if err != nil || len(block.(string)) == 0 {
		return line, err
	}

	check := rxRemoveEsc.ReplaceAll([]byte(line), []byte{' '})

	fork := lang.ShellProcess.Fork(lang.F_FUNCTION | lang.F_BACKGROUND | lang.F_CREATE_STDIN | lang.F_CREATE_STDOUT | lang.F_CREATE_STDERR)
	fork.Name.Set("(spellcheck)")
	fork.Stdin.SetDataType(types.Generic)
	_, err = fork.Stdin.Writeln(check)
	if err != nil {
		return line, err
	}

	_, err = fork.Execute([]rune(block.(string)))
	if err != nil {
		return line, err
	}

	b, err := fork.Stderr.ReadAll()
	if err != nil {
		return line, err
	}
	if len(b) != 0 {
		return line, fmt.Errorf("`config get shell spellcheck-func` STDERR: %s", string(utils.CrLfTrim(b)))
	}

	err = fork.Stdout.ReadArray(func(bWord []byte) {
		if len(bWord) == 0 {
			return
		}

		sWord := string(bytes.TrimSpace(bWord))

		if autocomplete.GlobalExes[sWord] || lang.MxFunctions.Exists(sWord) || lang.GoFunctions[sWord] != nil || lang.GlobalAliases.Exists(sWord) {
			return
		}

		if lang.ShellProcess.Variables.GetValue(sWord) != nil {
			return
		}

		if userdictionary.IsInDictionary(sWord) {
			return
		}

		line = strings.Replace(line, sWord, ansi.ExpandConsts("{UNDERLINE}"+sWord+"{UNDEROFF}"), -1)
	})

	return line, err
}
