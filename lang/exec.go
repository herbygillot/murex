package lang

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/lmorg/murex/builtins/pipes/null"
	"github.com/lmorg/murex/debug"
	"github.com/lmorg/murex/lang/types"
)

// External executes an external process.
func External(p *Process) error {
	if err := execute(p); err != nil {
		_, cmd := p.Exec.Get()
		if cmd != nil {
			p.ExitNum = cmd.ProcessState.ExitCode()
		} else {
			p.ExitNum = 1
		}
		return err

	}
	return nil
}

func execute(p *Process) error {
	p.Stdout.SetDataType(types.Generic)

	exeName, parameters, err := getCmdTokens(p)
	if err != nil {
		return err
	}
	cmd := exec.Command(exeName, parameters...)
	//cmd.Env = p.Exec.Env

	if p.HasCancelled() {
		return nil
	}

	//ctxCancel := p.Kill
	p.Kill = func() {
		if !debug.Enabled {
			defer func() { recover() }() // I don't care about errors in this instance since we are just killing the proc anyway
		}

		//ctxCancel()
		err := cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			name, _ := p.Args()
			os.Stderr.WriteString(
				fmt.Sprintf("\nError sending SIGTERM to `%s`: %s\n", name, err.Error()))
		}
	}

	switch {
	case p.IsMethod:
		cmd.Stdin = p.Stdin
	case p.Background.Get():
		cmd.Stdin = new(null.Null)
	default:
		cmd.Stdin = os.Stdin
	}

	if p.Stdout.IsTTY() {
		// If Stdout is a TTY then set the appropriate syscalls to allow the calling program to own the TTY....
		osSyscalls(cmd)
		cmd.Stdout = os.Stdout
	} else {
		// ....otherwise we just treat the program as a regular piped util
		cmd.Stdout = p.Stdout
	}

	// Pipe STDERR irrespective of whether the exec process is execting a TTY or not.
	// The reason for this is so that we can do some post-processing on the error stream (namely add colour to it),
	// however this might cause some bugs. If so please raise on github: https://github.com/lmorg/murex
	// In the meantime, you can force exec processes to write STDERR to the TTY via the `config` command in the shell:
	//
	//     config set proc force-tty true
	if p.Stderr.IsTTY() && forceTTY(p) {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = p.Stderr
	}

	if err := cmd.Start(); err != nil {
		//if !strings.HasPrefix(err.Error(), "signal:") {
		return err
		//}
	}

	p.Exec.Set(cmd.Process.Pid, cmd)

	if err := cmd.Wait(); err != nil {
		if !strings.HasPrefix(err.Error(), "signal:") {
			return err
		}
	}

	//debug.Log("exec env:",cmd.Env)

	return nil
}

func forceTTY(p *Process) bool {
	v, err := p.Config.Get("proc", "force-tty", types.Boolean)
	if err != nil {
		return false
	}
	return v.(bool)
}
