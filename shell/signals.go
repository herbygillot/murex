package shell

import (
	"fmt"
	"os"
	"regexp"

	"github.com/lmorg/murex/lang"
	"github.com/lmorg/murex/utils"
)

const (
	// PromptSIGINT defines the string to write when ctrl+c is pressed
	PromptSIGINT = "^C"

	// PromptSIGQUIT defines the string to write when ctrl+\ is pressed
	PromptSIGQUIT = "^\\"

	// PromptEOF defines the string to write when ctrl+d is pressed
	PromptEOF = "^D"
)

func sigint(interactive bool) {
	//os.Stderr.WriteString(PromptSIGINT)
	sigterm(interactive)
}

/*func sigtstp() {
	// see signals_unix.go
}*/

func sigterm(interactive bool) {
	if !interactive {
		os.Exit(0)
	}

	p := lang.ForegroundProc.Get()
	//p.Json("p =", p)

	switch {
	case p == nil:
		lang.ShellProcess.Stderr.Writeln([]byte("!!! Unable to identify forground process !!!"))
	case p.Kill == nil:
		lang.ShellProcess.Stderr.Writeln([]byte("!!! Unable to identify forground kill function !!!"))
	default:
		p.Kill()
	}
}

var rxWhiteSpace = regexp.MustCompilePOSIX(`[\r\n\t ]+`)

func sigquit(interactive bool) {
	if !interactive {
		os.Stderr.WriteString("Murex received SIGQUIT!" + utils.NewLineString)
		os.Exit(2)
	}

	//os.Stderr.WriteString(PromptSIGQUIT)
	os.Stderr.WriteString("Murex received SIGQUIT!" + utils.NewLineString)

	fids := lang.GlobalFIDs.ListAll()
	for _, p := range fids {
		if p.Kill != nil /*&& !p.HasTerminated()*/ {
			procName := p.Name.String()
			procParam, _ := p.Parameters.String(0)
			if procName == "exec" {
				procName = procParam
				procParam, _ = p.Parameters.String(1)
			}
			if len(procParam) > 60 {
				procParam = procParam[:60] + "..."
			}
			procParam = rxWhiteSpace.ReplaceAllString(procParam, " ")

			lang.ShellProcess.Stderr.Writeln([]byte(
				fmt.Sprintf(
					"!!! Force closing FID %d: %s %s !!!",
					p.Id, procName, procParam)))
			p.Kill()

			i, cmd := p.Exec.Get()
			if cmd != nil {
				err := cmd.Process.Kill()
				if err != nil {
					lang.ShellProcess.Stderr.Writeln([]byte(
						fmt.Sprintf(
							"!!! Error terminating FID %d (%d), `%s`: %s !!!",
							p.Id, i, procName, err.Error())))
				}
			}
		}
	}

	lang.ShellProcess.Stderr.Writeln([]byte("!!! Starting new prompt !!!"))
	go ShowPrompt()
}
