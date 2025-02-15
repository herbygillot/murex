# _murex_ Shell Docs

## Command Reference: `formap`

> Iterate through a map or other collection of data

## Description

`formap` is a generic tool for iterating through a map, table or other
sequences of data similarly like a `foreach`. In fact `formap` can even be
used on array too.

Unlike `foreach`, `formap`'s default output is `str`, so each new line will be
treated as a list item. This behaviour will differ if any additional flags are
used with `foreach`, such as `--jmap`.

## Usage

`formap` writes a list:

    <stdin> -> foreach variable { code-block } -> <stdout>
    
`formap` writes to a buffered JSON map:

    <stdin> -> formap --jmap key value { code-block (map key) } { code-block (map value) } -> <stdout>

## Examples

First of all lets assume the following dataset:

    set json people={
        "Tom": {
            "Age": 32,
            "Gender": "Male"
        },
        "Dick": {
            "Age": 43,
            "Gender": "Male"
        },
        "Sally": {
            "Age": 54,
            "Gender": "Female"
        }
    }
    
We can create human output from this:

    » $people -> formap key value { out "$key is $value[Age] years old" }
    Sally is 54 years old
    Tom is 32 years old
    Dick is 43 years old
    
> Please note that maps are intentionally unsorted so you cannot guarantee the
> order of the output produced even if the input has been superficially set in
> a specific order.

With `--jmap` we can turn that structure into a new structure:

    » $people -> formap --jmap key value { $key } { $value[Age] }
    {
        "Dick": "43",
        "Sally": "54",
        "Tom": "32"
    } 

## Flags

* `--jmap`
    Write a `json` map to STDOUT instead of an array

## Detail

`formap` can also work against arrays and tables as well. However `foreach` is
a much better tool for ordered lists and tables can look a little funky when
when there are more than 2 columns. In those instances you're better off using
`[` (index) to specify columns and then `tabulate` for any data transformation.

## See Also

* [commands/`[` (index)](../commands/index.md):
  Outputs an element from an array, map or table
* [commands/`for`](../commands/for.md):
  A more familiar iteration loop to existing developers
* [commands/`foreach`](../commands/foreach.md):
  Iterate through an array
* [types/`json` ](../types/json.md):
  JavaScript Object Notation (JSON) (primitive)
* [commands/`set`](../commands/set.md):
  Define a local variable and set it's value
* [commands/`tabulate`](../commands/tabulate.md):
  Table transformation tools
* [commands/`while`](../commands/while.md):
  Loop until condition false