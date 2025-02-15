//go:build !js
// +build !js

package main

import (
	"os"

	_ "github.com/lmorg/murex/builtins"
	"github.com/lmorg/murex/config/defaults"
	"github.com/lmorg/murex/config/profile"
	"github.com/lmorg/murex/debug"
	"github.com/lmorg/murex/lang"
	"github.com/lmorg/murex/shell"
	"github.com/lmorg/murex/utils/readline"
)

const (
	interactive    bool = true
	nonInteractive bool = false
)

func main() {
	readFlags()

	switch {
	case fRunTests:
		runTests()

	case fCommand != "":
		runCommandLine(fCommand)

	case len(fSource) > 0:
		runSource(fSource[0])

	default:
		startMurex()
	}

	debug.Log("[FIN]")
}

func runTests() error {
	lang.InitEnv()

	defaults.Defaults(lang.ShellProcess.Config, nonInteractive)
	shell.SignalHandler(nonInteractive)

	// compiled profile
	defaultProfile()

	// enable tests
	if err := lang.ShellProcess.Config.Set("test", "enabled", true, nil); err != nil {
		return err
	}
	if err := lang.ShellProcess.Config.Set("test", "auto-report", false, nil); err != nil {
		return err
	}
	if err := lang.ShellProcess.Config.Set("test", "verbose", false, nil); err != nil {
		return err
	}
	tty := readline.IsTerminal(int(os.Stdout.Fd()))
	if err := lang.ShellProcess.Config.Set("shell", "color", tty, nil); err != nil {
		return err
	}

	// run unit tests
	passed := lang.GlobalUnitTests.Run(lang.ShellProcess, "*")
	lang.ShellProcess.Tests.WriteResults(lang.ShellProcess.Config, lang.ShellProcess.Stdout)

	if !passed {
		os.Exit(1)
	}

	return nil
}

func runCommandLine(commandLine string) {
	lang.InitEnv()

	// default config
	defaults.Defaults(lang.ShellProcess.Config, nonInteractive)
	shell.SignalHandler(nonInteractive)

	// compiled profile
	defaultProfile()

	// load modules and profile
	if fLoadMods {
		profile.Execute(profile.F_PRELOAD | profile.F_MODULES | profile.F_PROFILE)
	}

	// read block from command line parameters
	execSource([]rune(commandLine), nil)
}

func runSource(filename string) {
	lang.InitEnv()

	// default config
	defaults.Defaults(lang.ShellProcess.Config, nonInteractive)
	shell.SignalHandler(nonInteractive)

	// compiled profile
	defaultProfile()

	// load modules a profile
	if fLoadMods {
		profile.Execute(profile.F_PRELOAD | profile.F_MODULES | profile.F_PROFILE)
	}

	// read block from disk
	disk, err := diskSource(filename)
	if err != nil {
		_, err := os.Stderr.WriteString(err.Error() + "\n")
		if err != nil {
			// wouldn't really make any difference at this point because we
			// cannot write to stderr anyway :(
			panic(err)
		}
		os.Exit(1)
	}
	execSource([]rune(string(disk)), nil)
}

func startMurex() {
	lang.InitEnv()

	// default config
	defaults.Defaults(lang.ShellProcess.Config, interactive)

	// compiled profile
	defaultProfile()

	// load modules and profile
	profile.Execute(profile.F_PRELOAD | profile.F_MODULES | profile.F_PROFILE)

	// start interactive shell
	shell.Start()
}
