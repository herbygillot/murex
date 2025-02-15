# _murex_ Shell Docs

## Parser Reference: And (`&&`) Logical Operator

> Continues next operation if previous operation passes

## Description

When in the **normal** run mode (see "schedulers" link below) this will only
run the command on the right hand side if the command on the left hand side
does not error. Neither STDOUT nor STDERR are piped.

This has no effect in `try` nor `trypipe` run modes because they automatically
apply stricter error handling.

## Examples

Second command runs because the first command doesn't error:

    » out: one && out: two
    one
    two
    
Second command does not run because the first command produces an error:

    » err: one && out: two
    one

## Detail

This is equivelent to a `try` block:

    try {
        err: one
        out: two
    }

## See Also

* [user-guide/Pipeline](../user-guide/pipeline.md):
  Overview of what a "pipeline" is
* [parser/STDERR Pipe (`?`) Token](../parser/pipe-err.md):
  Pipes STDERR from the left hand command to STDIN of the right hand command
* [user-guide/Schedulers](../user-guide/schedulers.md):
  Overview of the different schedulers (or 'run modes') in _murex_
* [commands/`err`](../commands/err.md):
  Print a line to the STDERR
* [commands/`out`](../commands/out.md):
  Print a string to the STDOUT with a trailing new line character
* [commands/`try`](../commands/try.md):
  Handles errors inside a block of code
* [commands/`trypipe`](../commands/trypipe.md):
  Checks state of each function in a pipeline and exits block on error