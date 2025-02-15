### Variables

There are two ways you can use variables with the math functions. Either by
string interpolation like you would normally with any other function, or
directly by name.

String interpolation:

```
» set abc=123
» = $abc==123
true
```

Directly by name:

```
» set abc=123
» = abc==123
false
```

To understand the difference between the two, you must first understand how
string interpolation works; which is where the parser tokenised the parameters
like so

```
command line: = $abc==123
token 1: command (name: "=")
token 2: parameter 1, string (content: "")
token 3: parameter 1, variable (name: "abc")
token 4: parameter 1, string (content: "==123")
```

Then when the command line gets executed, the parameters are compiled on demand
similarly to this crude pseudo-code

```
command: "="
parameters 1: concatenate("", GetValue(abc), "==123")
output: "=" "123==123"
```

Thus the actual command getting run is literally `123==123` due to the variable
being replace **before** the command executes.

Whereas when you call the variable by name it's up to `=` or `let` to do the
variable substitution.

```
command line: = abc==123
token 1: command (name: "=")
token 2: parameter 1, string (content: "abc==123")
```

```
command: "="
parameters 1: concatenate("abc==123")
output: "=" "abc==123"
```

The main advantage (or disadvantage, depending on your perspective) of using
variables this way is that their data-type is preserved.

```
» set str abc=123
» = abc==123
false

» set int abc=123
» = abc==123
true
```

Unfortunately is one of the biggest areas in _murex_ where you'd need to be
careful. The simple addition or omission of the dollar prefix, `$`, can change
the behavior of `=` and `let`.

### Strings

Because the usual _murex_ tools for encapsulating a string (`"`, `'` and `()`)
are interpreted by the shell language parser, it means we need a new token for
handling strings inside `=` and `let`. This is where backtick comes to our
rescue.

```
» set str abc=123
» = abc==`123`
true
```

Please be mindful that if you use string interpolation then you will need to
instruct `=` and `let` that your field is a string

```
» set str abc=123
» = `$abc`==`123`
true
```

### Best practice recommendation

As you can see from the sections above, string interpolation offers us some
conveniences when comparing variables of differing data-types, such as a `str`
type with a number (eg `num` or `int`). However it makes for less readable code
when just comparing strings. Thus the recommendation is to avoid using string
interpolation except only where it really makes sense (ie use it sparingly).

### Non-boolean logic

Thus far the examples given have been focused on comparisons however `=` and
`let` supports all the usual arithmetic operators:

```
» = 10+10
20

» = 10/10
1

» = (4 * (3 + 2))
20

» = `foo`+`bar`
foobar
```

### Read more

_murex_ uses the [govaluate package](https://github.com/Knetic/govaluate). More information can be found in it's manual.