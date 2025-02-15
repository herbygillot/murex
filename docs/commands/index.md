# _murex_ Shell Docs

## Command Reference: `[` (index)

> Outputs an element from an array, map or table

## Description

Outputs an element or multiple elements from an array, map or table.

Please note that indexes in _murex_ are counted from zero.

## Usage

    <stdin> -> [ element ] -> <stdout>
    $variable[ element ] -> <stdout>
    
    <stdin> -> ![ element ] -> <stdout>

## Examples

Return the 2nd (1), 4th (3) and 6th (5) element in an array

    » ja [0..9] -> [ 1 3 5 ]
    [
        "1",
        "3",
        "5"
    ]
    
Return the data-type and description of **config shell syntax-highlighting**

    » config -> [[ /shell/syntax-highlighting ]] -> [ Data-Type Description ]
    [
        "bool",
        "Syntax highlighting of murex code when in the interactive shell"
    ]
    
Return all elements _except_ for 1 (2nd), 3 (4th) and 5 (6th)

    » a: [0..9]-> ![ 1 3 5 ]
    0
    2
    4
    6
    7
    8
    9
    
Return all elements except for the data-type and description

    » config -> [[ /shell/syntax-highlighting ]] -> ![ Data-Type Description ]
    {
        "Default": true,
        "Dynamic": false,
        "Global": true,
        "Value": true
    }

## Detail

### Index counts from zero

Indexes in _murex_ behave like any other computer array in that all arrays
start from zero (`0`).

### Include vs exclude

As demonstrated in the examples above, `[` specifies elements to include
where as `![` specifies elements to exclude.

### Don't error upon missing elements

By default, **index** generates an error if an element doesn't exist. However
you can disable this behavior in `config`

    » config -> [ foobar ]
    Error in `[` ((builtin) 2,11): Key 'foobar' not found
    
    » config set index silent true
    
    » config -> [ foobar ]

## Synonyms

* `[`
* `![`


## See Also

* [commands/`@[` (range) ](../commands/range.md):
  Outputs a ranged subset of data from STDIN
* [commands/`[[` (element)](../commands/element.md):
  Outputs an element from a nested structure
* [commands/`a` (mkarray)](../commands/a.md):
  A sophisticated yet simple way to build an array or list
* [commands/`config`](../commands/config.md):
  Query or define _murex_ runtime settings
* [commands/`count`](../commands/count.md):
  Count items in a map, list or array
* [commands/`ja` (mkarray)](../commands/ja.md):
  A sophisticated yet simply way to build a JSON array
* [commands/`mtac`](../commands/mtac.md):
  Reverse the order of an array