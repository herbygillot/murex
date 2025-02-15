# _murex_ Shell Docs

## Command Reference: `autocomplete`

> Set definitions for tab-completion in the command line

## Description

`autocomplete` digests a JSON schema and uses that to define the tab-
completion rules for suggestions in the interactive command line.

## Usage

    autocomplete get [ command ] -> <stdout>
    
    autocomplete set command { mxjson }

## Flags

* `get`
    output all autocompletion schemas
* `set`
    define a new autocompletion schema

## Detail

### Undefining autocomplete

Currently there is no support for undefining an autocompletion rule however
you can overwrite existing rules.

## Directives

The directives are listed below. Headings are formatted as follows:

    "DirectiveName": json data-type (default value)
    
Where "default value" is what will be auto-populated at run time if you
don't define an autocomplete schema manually.

### "Alias": string ("")

Aliases are used inside **FlagValues** as a way of pointing one flag to another
without duplicating code. eg `-v` and `--version` might be the same flag. Or
`-?`, `-h` and `--help`. With **Alias** you can write the definitions for one
flag and then point all the synonyms as an alias to that definition.

### "AllowAny": boolean (false)

The way autocompletion works in _murex_ is the suggestion engine looks for
matches and if it fines one, it then moves onto the next index in the JSON
schema. This means unexpected values typed in the interactive terminal will
break the suggestion engine's ability to predict what the next expected
parameter should be. Setting **AllowAny** to `true` tells the suggestion
engine to accept any value as the next parameter thus allowing it to then
predict the next parameter afterwards.

This directive isn't usually necessary because such fields are often the last
parameter or most parameters can be detectable with a reasonable amount of
effort. However **AllowAny** is often required for more complex command line
tools.

### "AllowMultiple": boolean (false)

Set to `true` to enable multiple parameters following the same rules as defined
in this index. For example the following will suggest directories on each tab
for multiple parameters:

    autocomplete set example { [{
        "IncDirs": true,
        "AllowMultiple": true
    }] }
    
### "AnyValue": boolean (false)

Deprecated. Please use **AllowAny** instead.

### "AutoBranch": boolean (false)

Use this in conjunction with **Dynamic**. If the return is an array of paths,
for example `[ "/home/foo", "/home/bar" ]` then **AutoBranch** will return
the following patterns in the command line:

    » example [tab]
    # suggests "/home/"
    
    » example /home/[tab]
    # suggests "/home/foo" and "/home/bar"
    
Please note that **AutoBranch**'s behavior is also dependant on a "shell"
`config` setting, recursive-enabled":

    » config get shell recursive-enabled
    true
    
### "Dynamic": string ("")

This is a _murex_ block which returns an array of suggestions.

Code inside that block are executed like a function and the parameters will
mirror the same as those parameters entered in the interactive terminal.

### "DynamicDesc": string ("")

This is very similar to **Dynamic** except your function should return a
map instead of an array. Where each key is the suggestion and the value is
a description.

The description will appear either in the hint text or alongside the
suggestion - depending on which suggestion "popup" you define (see
**ListView**).

### "ExecCmdline": boolean (false)

Sometimes you'd want your autocomplete suggestions to aware of the output
returned from the commands that preceded it. For example the suggestions
for `[` (index) will depend entirely on what data is piped into it.

**ExecCmdline** tells _murex_ to run the commandline up until the command
which your cursor is editing and pipe that output to the STDIN of that
commands **Dynamic** or **DynamicDesc** code block.

> This is a dangerous feature to enable so **ExecCmdline** is only honoured
> if the commandline is considered "safe". **Dynamic** / **DynamicDesc**
> will still be executed however if the commandline is "unsafe" then your
> dynamic autocompletion blocks will have no STDIN.

Because this is a dangerous feature, your partial commandline will only
execute if the following conditions are met:

* the commandline must be one pipeline (eg `;` tokens are not allowed)
* the commandline must not have any new line characters
* there must not be any redirection, including named pipes
  (eg `cmd <namedpipe>`) and the STDOUT/STDERR switch token (`?`)
* the commandline doesn't inline any variables (`$strings`, `@arrays`) or
  functions (`${subshell}`, `$[index]`)
* lastly all commands are whitelisted in "safe-commands"
  (`config get shell safe-commands`)

If these criteria are met, the commandline is considered "safe"; if any of
those conditions fail then the commandline is considered "unsafe".

_murex_ will come with a number of sane commands already included in its
`safe-commands` whitelist however you can add or remove them using `config`

    » function: foobar { -> match foobar }
    » config: eval shell safe-commands { -> append foobar }
    
Remember that **ExecCmdline** is designed to be included with either
**Dynamic** or **DynamicDesc** and those code blocks would need to read
from STDIN:

    autocomplete set "[" { [{
        "AnyValue": true,
        "AllowMultiple": true,
        "ExecCmdline": true,
        "Dynamic": ({
            switch ${ get-type: stdin } {
                case * {
                    <stdin> -> [ 0: ] -> format json -> [ 0 ]
                }
                
                catch {
                    <stdin> -> formap k v { out $k } -> cast str -> append "]"
                }
            }
        })
    }] }
    
### "FileRegexp": string ("")

When set in conjunction with **IncFiles**, this directive will filter on files
files which match the regexp string. eg to only show ".txt" extensions you can
use the following:

    autocomplete set notepad.exe { [{
        "IncFiles": true,
        "FileRegexp": (\.txt)
    }] }
    
> Please note that you may need to double escape any regexp strings: escaping
> the `.` match and then also escaping the escape character in JSON. It is
> recommended you use the `mxjson` method of quoting using parentheses as this
> will compile that string into JSON, automatically adding additional escaping
> where required.

### "FlagValues": map of arrays (null)

This is a map of the flags with the values being the same array of directive
as the top level.

This allows you to nest operations by flags. eg when a flag might accept
multiple parameters.

**FlagValues** takes a map of arrays, eg

    autocomplete set example { [{
        "Flags": [ "add", "delete" ],
        "FlagValues": {
            "add": [{
                "Flags": [ "foo" ]
            }],
            "delete": [{
                "Flags": [ "bar" ]
            }]
        }
    }] }
    
...will provide "foo" as a suggestion to `example add`, and "bar" as a
suggestion to `example delete`.

#### Defaults for matched flags

You can set default properties to all matched flags by using `*` as a
**FlagValues** value. To expand the above example...

    autocomplete set example { [{
        "Flags": [ "add", "delete" ],
        "FlagValues": {
            "add": [{
                "Flags": [ "foo" ]
            }],
            "delete": [{
                "Flags": [ "bar" ]
            }],
            "*": [{
                "IncFiles"
            }]
        }
    }] }
    
...in this code we are saying not only does "add" support "foo" and "delete"
supports "bar", but both "add" and "delete" also supports any filesystem files.

This default only applies if there is a matched **Flags** or **FlagValues**.

#### Defaults for any flags (including unmatched)

If you wanted a default which applied to all **FlagValues**, even when the flag
wasn't matched, then you can use a zero length string (""). For example

    autocomplete set example { [{
        "Flags": [ "add", "delete" ],
        "FlagValues": {
            "add": [{
                "Flags": [ "foo" ]
            }],
            "delete": [{
                "Flags": [ "bar" ]
            }],
            "": [{
                "IncFiles"
            }]
        }
    }] }
    
### "Flags": array of strings (auto-populated from man pages)

Setting **Flags** is the fastest and easiest way to populate suggestions
because it is just an array of strings. eg

    autocomplete set example { [{
        "Flags": [ "foo", "bar" ]
    }] }
    
If a command doesn't **Flags** already defined when you request a completion
suggestion but that command does have a man page, then **Flags** will be
automatically populated with any flags identified from an a quick parse of
the man page. However because man pages are written to be human readable
rather than machine parsable, there may not be a 100% success rate with the
automatic man page parsing.
    
### "FlagsDesc": map of strings (null)

This is the same concept as **Flags** except it is a map with the suggestion
as a key and description as a value. This distinction is the same as the
difference between **Dynamic** and **DynamicDesc**.

Please note that currently man page parsing cannot provide a description so
only **Flags** get auto-populated.

### "Goto": string ("")

This is a `goto` in programming terms. While "ugly" it does allow for quick and
easy structural definitions without resorting to writing the entire
autocomplete in code.

**Goto** takes a string which represents the path to jump to from the top level
of that autocomplete definition. The path should look something like:
`/int/string/int/string....` where

* the first character is the separator,

* the first value is an integer that relates to the index in your autocomplete
  array,

* the second value is a string which points to the flag value map (if you
  defined **FlagValues**),

* the third value is the integer of the autocomplete array inside that
**FlagValues** map,

* ...and so on as necessary.

An example of a really simple **Goto**:

    autocomplete set dd { [
        {
            "Flags": [ "if=", "of=", "bs=", "iflag=", "oflag=", "count=", "status=" ],
            "FlagValues": {
                "if": [{ 
                    "IncFiles": true
                }],
                "of": [{ 
                    "IncFiles": true
                }],
                "*": [{
                    "AllowAny": true
                }]
            }
        },
        {
            "Goto": "/0"
        }
    ] }
    
**Goto** is given precedence over any other directive. So ensure it's the only
directive in it's group.

### "IncDirs": boolean (false)

Enable to include directories.

Not needed if **IncFiles** is set to `true`.

Behavior of this directive can be altered with `config set shell
recursive-enabled`

### "IncExeAll": boolean (false)

Enable this to any executables. Suggestions will include aliases, functions
builtins and any executables in `$PATH`. It will not include private functions.

### "IncExePath": boolean (false)

Enable this to include any executables in `$PATH`. Suggestions will not include
aliases, functions nor privates.

### "IncFiles": boolean (true)

Include files and directories. This is enabled by default for any commands
that don't have autocomplete defined but you will need to manually enable
it in any `autocomplete` schemas you create and want files as part of the
suggestions.

If you want to filter files based on file name then you can set a regexp
string to match to using **FileRegexp**.

### "IncManPage": boolean (false)

The default behavior for commands with no autocomplete defined is to parse the
man page and use those results. If a custom autocomplete is defined then that
man page parser is disabled by default. You can re-enable it and include its
results with other flags and behaviors you define by using this directive.

### "ListView": boolean (false)

This alters the appearance of the autocompletion suggestions "popup". Rather
than suggestions being in a grid layout (with descriptions overwriting the
hint text) the suggestions are in a list view with the descriptions next to
them on the same row (similar to how an IDE might display it's suggestions).

### "NestedCommand": boolean (false)

Only enable this if the command you are autocompleting is a nested parameter
of the parent command you have types. For example with `sudo`, once you've
typed the command name you wish to elivate, then you would want suggestions
for that command rather than for `sudo` itself.

### "Optional": boolean (false)

Specifies if a match is required for the index in this schema. ie optional
flags.

## See Also

* [commands/`<stdin>` ](../commands/stdin.md):
  Read the STDIN belonging to the parent code block
* [commands/`[` (index)](../commands/index.md):
  Outputs an element from an array, map or table
* [commands/`alias`](../commands/alias.md):
  Create an alias for a command
* [commands/`config`](../commands/config.md):
  Query or define _murex_ runtime settings
* [commands/`function`](../commands/function.md):
  Define a function block
* [commands/`get-type`](../commands/get-type.md):
  Returns the data-type of a variable or pipe
* [commands/`private`](../commands/private.md):
  Define a private function block
* [commands/`summary` ](../commands/summary.md):
  Defines a summary help text for a command
* [commands/`switch`](../commands/switch.md):
  Blocks of cascading conditionals
* [types/mxjson](../types/mxjson.md):
  Murex-flavoured JSON (primitive)