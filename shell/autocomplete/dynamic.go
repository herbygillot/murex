package autocomplete

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lmorg/murex/builtins/pipes/streams"
	"github.com/lmorg/murex/debug"
	"github.com/lmorg/murex/lang"
	"github.com/lmorg/murex/lang/types"
	"github.com/lmorg/murex/utils"
	"github.com/lmorg/murex/utils/lists"
	"github.com/lmorg/murex/utils/parser"
)

type dynamicArgs struct {
	exe    string
	params []string
	float  int
}

func matchDynamic(f *Flags, partial string, args dynamicArgs, act *AutoCompleteT) {
	// Default to building up from Dynamic field. Fall back to DynamicDefs
	dynamic := f.Dynamic
	if f.Dynamic == "" {
		dynamic = f.DynamicDesc
	}
	if dynamic == "" {
		return
	}

	if !types.IsBlock([]byte(dynamic)) {
		lang.ShellProcess.Stderr.Writeln([]byte("Dynamic autocompleter is not a code block"))
		return
	}
	block := []rune(dynamic[1 : len(dynamic)-1])

	softTimeout, _ := lang.ShellProcess.Config.Get("shell", "autocomplete-soft-timeout", types.Integer)
	//if err != nil {
	//	softTimeout = 100
	//}

	hardTimeout, _ := lang.ShellProcess.Config.Get("shell", "autocomplete-hard-timeout", types.Integer)
	//if err != nil {
	//	hardTimeout = 5000
	//}

	softCtx, _ := context.WithTimeout(context.Background(), time.Duration(int64(softTimeout.(int)))*time.Millisecond)
	hardCtx, _ := context.WithTimeout(context.Background(), time.Duration(int64(hardTimeout.(int)))*time.Millisecond)
	wait := make(chan bool)
	done := make(chan bool)

	act.largeMin()
	/*if f.ListView {
		// check this here so delayed results can still be ListView
		// (ie after &act has timed out)
		act.TabDisplayType = readline.TabDisplayList
	}*/

	go func() {
		// Run the commandline if ExecCmdline flag set AND commandline considered safe
		var fStdin int
		cmdlineStdout := streams.NewStdin()
		if f.ExecCmdline && !act.ParsedTokens.Unsafe {
			cmdline := lang.ShellProcess.Fork(lang.F_BACKGROUND | lang.F_NO_STDIN | lang.F_NO_STDERR)
			cmdline.Stdout = cmdlineStdout
			cmdline.Name.Set(args.exe)
			cmdline.FileRef = ExesFlagsFileRef[args.exe]
			cmdline.Execute(act.ParsedTokens.Source[:act.ParsedTokens.LastFlowToken])

		} else {
			fStdin = lang.F_NO_STDIN
		}

		// don't share incomplete parameters with dynamic autocompletion blocks
		params := act.ParsedTokens.Parameters
		switch len(params) {
		case 0: // do nothing
		case 1:
			params = []string{}
		default:
			params = params[:len(params)-1]
		}

		// Execute the dynamic code block
		fork := lang.ShellProcess.Fork(lang.F_FUNCTION | lang.F_NEW_MODULE | lang.F_BACKGROUND | fStdin | lang.F_CREATE_STDOUT | lang.F_NO_STDERR)
		fork.Name.Set(args.exe)
		fork.Parameters.DefineParsed(params)
		fork.FileRef = ExesFlagsFileRef[args.exe]
		if f.ExecCmdline && !act.ParsedTokens.Unsafe {
			fork.Stdin = cmdlineStdout
		}
		fork.Process.Variables.Set(fork.Process, "ISMETHOD", act.ParsedTokens.PipeToken != parser.PipeTokenNone, types.Boolean)
		exitNum, err := fork.Execute(block)

		if err != nil {
			lang.ShellProcess.Stderr.Writeln([]byte("dynamic autocomplete code could not compile: " + err.Error()))
		}
		if exitNum != 0 && debug.Enabled {
			lang.ShellProcess.Stderr.Writeln([]byte("dynamic autocomplete returned a none zero exit number." + utils.NewLineString))
		}

		select {
		case <-hardCtx.Done():
			act.ErrCallback(fmt.Errorf("dynamic autocompletion took too long"))
			return
		default:
		}

		if f.Dynamic != "" {
			var (
				timeout bool
				items   []string
			)

			select {
			case <-softCtx.Done():
				timeout = true
			default:
				wait <- true
			}

			var (
				incFiles    bool
				incDirs     bool
				incExePath  bool
				incExeAll   bool
				incManPages bool
			)

			err := fork.Stdout.ReadArray(func(b []byte) {
				//s := string(bytes.TrimSpace(b))
				s := string(b)

				if len(s) == 0 {
					return
				}

				switch s {
				case "@IncFiles":
					incFiles = true
				case "@IncDirs":
					incDirs = true
				case "@IncExePath":
					incExePath = true
				case "@incExeAll":
					incExeAll = true
				case "@IncManPages":
					incManPages = true

				default:
					if strings.HasPrefix(s, partial) {
						items = append(items, s[len(partial):])
					}
				}
			})

			switch {
			case incFiles:
				files := matchFilesAndDirs(partial, act)
				items = append(items, files...)
			case incDirs:
				files := matchDirs(partial, act)
				items = append(items, files...)
			case incExePath:
				pathexes := allExecutables(false)
				items = append(items, matchExes(partial, pathexes, false)...)
			case incExeAll:
				pathall := allExecutables(true)
				items = append(items, matchExes(partial, pathall, false)...)
			case incManPages:
				flags := scanManPages(args.exe)
				items = append(items, lists.CropPartial(flags, partial)...)
			}

			if err != nil {
				debug.Log(err)
			}

			if f.AutoBranch && !act.CacheDynamic {
				autoBranch(&items)
			}

			if timeout {
				formatSuggestionsArray(act.ParsedTokens, items)
				act.DelayedTabContext.AppendSuggestions(items)
				//act.ErrCallback(fmt.Errorf(" ")) // clear slow autocomplete message
			} else {
				act.append(items...)
			}

		} else {
			var (
				timeout bool
				items   = make(map[string]string)
			)

			select {
			case <-softCtx.Done():
				timeout = true
			default:
				wait <- true
			}

			fork.Stdout.ReadMap(lang.ShellProcess.Config, func(key string, value string, last bool) {
				if strings.HasPrefix(key, partial) {
					value = strings.Replace(value, "\r", "", -1)
					value = strings.Replace(value, "\n", " ", -1)

					if timeout {
						items[key[len(partial):]] = value
					} else {
						act.appendDef(key[len(partial):], value)
					}
				}
			})

			if timeout {
				formatSuggestionsMap(act.ParsedTokens, &items)
				act.DelayedTabContext.AppendDescriptions(items)
			}
		}

		done <- true
		//act.ErrCallback(fmt.Errorf(" ")) // clear slow autocomplete message
	}()

	select {
	case <-done:
		//act.MinTabItemLength = 0
		return
	case <-wait:
		<-done
		//act.MinTabItemLength = 0
		return

	case <-softCtx.Done():
		if len(act.Items) == 0 && len(act.Definitions) == 0 {
			act.ErrCallback(fmt.Errorf("long running autocompletion pushed to the background"))
			//act.appendDef("", "")
			//act.MinTabItemLength = 0
		}

		return
	}

}
