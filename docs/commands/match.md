# _murex_ Shell Docs

## Command Reference: `match`

> Match an exact value in an array

## Description

`match` takes input from STDIN and returns any array items / lines which
contain an exact match of the parameters supplied.

When multiple parameters are supplied they are concatenated into the search
string and white space delimited. eg all three of the below are the same:

    match "a b c"
    match a\sb\sc
    match a b c
    match a    b    c
    
If you want to return everything except the search string then use `!match

## Usage

Match every occurrence of search string

    <stdin> -> match search string -> <stdout>
    
Match everything except search string

    <stdin> -> !match search string -> <stdout>

## Examples

Match **Wed**

    » ja: [Monday..Friday] -> match Wed
    [
        "Wednesday"
    ]
    
Match everything except **Wed**

    » ja: [Monday..Friday] -> !match Wed
    [
        "Monday",
        "Tuesday",
        "Thursday",
        "Friday"
    ] 

## Detail

`match` is data-type aware so will work against lists or arrays of whichever
_murex_ data-type is passed to it via STDIN and return the output in the
same data-type.

## Synonyms

* `match`
* `!match`


## See Also

* [commands/`2darray` ](../commands/2darray.md):
  Create a 2D JSON array from multiple input sources
* [commands/`a` (mkarray)](../commands/a.md):
  A sophisticated yet simple way to build an array or list
* [commands/`append`](../commands/append.md):
  Add data to the end of an array
* [commands/`count`](../commands/count.md):
  Count items in a map, list or array
* [commands/`ja` (mkarray)](../commands/ja.md):
  A sophisticated yet simply way to build a JSON array
* [commands/`jsplit` ](../commands/jsplit.md):
  Splits STDIN into a JSON array based on a regex parameter
* [commands/`map` ](../commands/map.md):
  Creates a map from two data sources
* [commands/`msort` ](../commands/msort.md):
  Sorts an array - data type agnostic
* [commands/`prefix`](../commands/prefix.md):
  Prefix a string to every item in a list
* [commands/`prepend` ](../commands/prepend.md):
  Add data to the start of an array
* [commands/`pretty`](../commands/pretty.md):
  Prettifies JSON to make it human readable
* [commands/`regexp`](../commands/regexp.md):
  Regexp tools for arrays / lists of strings
* [commands/`suffix`](../commands/suffix.md):
  Prefix a string to every item in a list
* [commands/`ta` (mkarray)](../commands/ta.md):
  A sophisticated yet simple way to build an array of a user defined data-type