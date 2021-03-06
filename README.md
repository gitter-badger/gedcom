github.com/elliotchance/gedcom
==============================

[![Build Status](https://travis-ci.org/elliotchance/gedcom.svg?branch=master)](https://travis-ci.org/elliotchance/gedcom)
[![GoDoc](https://godoc.org/github.com/elliotchance/gedcom?status.svg)](https://godoc.org/github.com/elliotchance/gedcom)
[![codecov](https://codecov.io/gh/elliotchance/gedcom/branch/master/graph/badge.svg)](https://codecov.io/gh/elliotchance/gedcom)

`gedcom` is an advanced Go-style library for encoding, decoding, traversing,
exporting and diffing GEDCOM files.

- [Project Goals](#project-goals)
- [Installation](#installation)
- [Command Line Tools](#command-line-tools)
  * [Comparing GEDCOM Files](#comparing-gedcom-files)
  * [Rendering as HTML](#rendering-as-html)
  * [Converting to JSON](#converting-to-json)
  * [Converting to Text](#converting-to-text)
  * [Using the Checksum File](#using-the-checksum-file)
- [Library](#library)
  * [Decoding a Document](#decoding-a-document)
  * [Encoding a Document](#encoding-a-document)
  * [Traversing a Document](#traversing-a-document)
  * [Comparing & Diffing Nodes](#comparing---diffing-nodes)
  * [Dates](#dates)
  * [Node Equality](#node-equality)
  * [Tree Walking & Filtering](#tree-walking---filtering)
  * [Individual Names](#individual-names)


Project Goals
=============

1. Support all GEDCOM files by supporting the encoding and not the GEDCOM
standard itself. Many GEDCOM libraries that try to follow the standard run into
trouble when applications do not follow the same standards or the standard is
interpreted differently. `gedcom` retains all tags and structure in the original
GEDCOM file.

2. Build structures and functions that provide a nicer API for common
operations. For example iterating through individuals in a file, traversing
through family connections, understanding dates, etc.

3. Export to other file formats. Such as HTML, JSON, text, etc. So that
information can be manipulated and ingested by other applications.

4. Provide more advanced functionality to deal with comparing and diffing GEDCOM
files.

Installation
============

You can download the latest binaries for Mac, Windows or Linux on the
[Releases][15] page. This will not require you to install Go or any other
dependencies.

If you wish to build it from source you must include the dependencies with:

```bash
dep ensure
```


Command Line Tools
==================

Comparing GEDCOM Files
----------------------

The `gedcomdiff` CLI tool will generate a HTML report when comparing two GEDCOM
files. There are lots of configurable options, but the most basic usage is:

```bash
gedcomdiff -left-gedcom file1.ged -right-gedcom file2.ged
```

Rendering as HTML
-----------------

`gedcom2html` converts a GEDCOM file to a directory of HTML files. This produces
a pretty output that looks like this:
[http://dechauncy.family](http://dechauncy.family)

```txt
Usage of gedcom2html:
  -gedcom string
    	Input GEDCOM file.
  -output-dir string
    	Output directory. It will use the current directory if output-dir is not provided. Output files will only be added or replaced. Existing files will not be deleted. (default ".")
```

Converting to JSON
------------------

`gedcom2json` is a subpackage and binary that converts a GEDCOM file to a JSON
structure. It offers several options for the output:

```
Usage of gedcom2json:
  -exclude-tags string
    	Comma-separated list of tags to ignore.
  -gedcom string
    	Input GEDCOM file.
  -no-pointers
    	Do not include Pointer values ("ptr" attribute) in the output JSON. This is useful to activate when comparing GEDCOM files that have had pointers generated from different sources.
  -only-official-tags
    	Only include tags from the GEDCOM standard in the output.
  -only-tags string
    	Only include these tags in the output.
  -pretty-json
    	Pretty print JSON.
  -pretty-tags
    	Output tags with their descriptive name instead of their raw tag value. For example, "BIRT" would be output as "Birth".
  -single-name
    	When there are multiple names for an individual this will return the first of the name nodes only.
  -string-name
    	Convert NAME tags to a string (instead of the object parts).
  -tag-keys
    	Use tags (pretty or raw) as object keys rather than arrays.
```

Converting to Text
------------------

`gedcom2text` is a subpackage and binary that converts a GEDCOM file to a simple
text output (or split into individual files) that is ideal for easily reading
(by a person) and designed to be as friendly as possible when using diff tools.

```
Usage of gedcom2text:
  -gedcom string
    	Input GEDCOM file.
  -no-change-times
    	Do not change timestamps.
  -no-empty-deaths
    	Do not include Death node if there are no visible details.
  -no-places
    	Do not include places.
  -no-sources
    	Do not include sources.
  -only-official-tags
    	Only output official GEDCOM tags.
  -single-name
    	Only output the primary name.
  -split-dir string
    	Split the individuals into separate files in this directory.
```

Using the Checksum File
-----------------------

The `-checksum` option for `gedcom2html` generates a file called `checksum.csv`.
This file contains the file name and SHA-1 checksum like:

```
amos-adams.html,b0538fb8186a50c4079c902fec2b4ba0af843061
massachusetts-united-states.html,79db811c089e8ab5653d34551e6540cb2ea2c947
```

The lines are ordered by the file name so the output is ideal for comparison.

Here is an example of using the previous and current checksum file to generate
sync commands:

```bash
join -a 1 -a 2 -t, -o 0.1,1.2,2.2 /old/checksum.csv /new/checksum.csv | \
    awk -F, '$2 == $3 { next } { print $3 == "" \
        ? "rm /some/folder/" $1 \
        : "cp" " " $1 " /some/folder/" $1 }'
```

Will produce commands like:

```
cp abos-adams.html /some/folder/abos-adams.html
rm /some/folder/massachusetts-united-states.html
```


Library
=======

Decoding a Document
-------------------

Decoding a GEDCOM stream:

```go
ged := "0 HEAD\n1 CHAR UTF-8"

decoder := gedcom.NewDecoder(strings.NewReader(ged))
document, err := decoder.Decode()
if err != nil {
    panic(err)
}
```

If you are reading from a file you can use `NewDocumentFromGEDCOMFile`:

```go
document, err := gedcom.NewDocumentFromGEDCOMFile("family.ged")
if err != nil {
    panic(err)
}
```

Encoding a Document
-------------------

```go
buf := bytes.NewBufferString("")

encoder := NewEncoder(buf, doc)
err := encoder.Encode()
if err != nil {
	panic(err)
}
```

If you need the GEDCOM data as a string you can simply using `fmt.Stringer`:

```go
data := document.String()
```

Traversing a Document
---------------------

On top of the raw document is a powerful API that takes care of the complex
traversing of the Document. Here is a simple example:

```go
for _, individual := range document.Individuals() {
    fmt.Println(individual.Name().String())
}
```

Some of the nodes in a GEDCOM file have been replaced with more function rich
types, such as names, dates, families and more. See
[godoc](https://godoc.org/github.com/elliotchance/gedcom) for a complete list of
API methods.


Comparing & Diffing Nodes
-------------------------

The [`CompareNodes`][1] recursively compares two nodes. For example:

```
0 INDI @P3@           |  0 INDI @P4@
1 NAME John /Smith/   |  1 NAME J. /Smith/
1 BIRT                |  1 BIRT
2 DATE 3 SEP 1943     |  2 DATE Abt. Sep 1943
1 DEAT                |  1 BIRT
2 PLAC England        |  2 DATE 3 SEP 1943
1 BIRT                |  1 DEAT
2 DATE Abt. Oct 1943  |  2 DATE Aft. 2001
                      |  2 PLAC Surry, England
```

Produces a [`*NodeDiff`][2] than can be rendered with the [`String`][3] method:

```
LR 0 INDI @P3@
L  1 NAME John /Smith/
LR 1 BIRT
L  2 DATE Abt. Oct 1943
LR 2 DATE 3 SEP 1943
 R 2 DATE Abt. Sep 1943
LR 1 DEAT
L  2 PLAC England
 R 2 DATE Aft. 2001
 R 2 PLAC Surry, England
 R 1 NAME J. /Smith/
```

Dates
-----

Dates in GEDCOM files can be very complex as they can cater for many scenarios:

1. Incomplete, like "Dec 1943"
2. Anchored, like "Aft. 3 Sep 2003" or "Before 1923"
3. Ranges, like "Bet. 4 Apr 1823 and 8 Apr 1823"

This package provides a very rich API for dealing with all kind of dates in a
meaningful and sensible way. Some notable features include:

1. All dates, even though that specify an specific day have a minimum and
maximum value that are their true bounds. This is especially important for
larger date ranges like the whole month of "Jun 1945".
2. Upper and lower bounds of dates can be converted to the native Go `time.Time`
object.
3. There is a `Years` function that provides a convenient way to normalise a
date range into a number for easier distance and comparison measurements.
4. Algorithms for calculating the similarity of dates on a configurable
parabola.

Node Equality
-------------

[`Node.Equals`][9] performs a shallow comparison between two nodes. The
implementation is different depending on the types of nodes being compared. You
should see the specific documentation for the Node.

Equality is not to be confused with the `Is` function seen on some of the nodes,
such as [`Date.Is`][12]. The `Is` function is used to compare exact raw values
in nodes.

[`DeepEqual`][10] tests if left and right are recursively equal.

Tree Walking & Filtering
------------------------

The [`Filter`][4] function recursively removes or manipulates nodes with a
[`FilterFunction`][5]:

```go
newNodes := gedcom.Filter(node, func (node gedcom.Node) (gedcom.Node, bool) {
    if node.Tag().Is(gedcom.TagIndividual) {
        // false means it will not traverse children, since an
        // individual can never be inside of another individual.
        return node, false
    }

    return nil, false
})

// Remove all tags that are not official.
newNodes := gedcom.Filter(node, gedcom.OfficialTagFilter())
```

Filter functions:

1. [`BlacklistTagFilter`][7]
2. [`OfficialTagFilter`][8]
3. [`SimpleNameFilter`][12]
4. [`WhitelistTagFilter`][6]

Individual Names
----------------

A [`NameNode`][14] represents all the parts that make up a single name. An
individual may have more than one name, each one would be represented by a
`NameNode`.

Apart from functions to extract name parts there is also [`Format`][13] which
works similar to `fmt.Printf` where placeholders represent different components
of the name:

```txt
%% "%"
%f GivenName
%l Surname
%m SurnamePrefix
%p Prefix
%s Suffix
%t Title
```

Each of the letters may be in upper case to convert the name part to upper case
also. Whitespace before, after and between name components will be removed:

```txt
name.Format("%l, %f")     // Smith, Bob
name.Format("%f %L")      // Bob SMITH
name.Format("%f %m (%l)") // Bob (Smith)
```


[1]: https://godoc.org/github.com/elliotchance/gedcom#CompareNodes
[2]: https://godoc.org/github.com/elliotchance/gedcom#NodeDiff
[3]: https://godoc.org/github.com/elliotchance/gedcom#NodeDiff.String
[4]: https://godoc.org/github.com/elliotchance/gedcom#Filter
[5]: https://godoc.org/github.com/elliotchance/gedcom#FilterFunction
[6]: https://godoc.org/github.com/elliotchance/gedcom#WhitelistTagFilter
[7]: https://godoc.org/github.com/elliotchance/gedcom#BlacklistTagFilter
[8]: https://godoc.org/github.com/elliotchance/gedcom#OfficialTagFilter
[9]: https://godoc.org/github.com/elliotchance/gedcom#Node
[10]: https://godoc.org/github.com/elliotchance/gedcom#DeepEqual
[11]: https://godoc.org/github.com/elliotchance/gedcom#Date.Is
[12]: https://godoc.org/github.com/elliotchance/gedcom#SimpleNameFilter
[13]: https://godoc.org/github.com/elliotchance/gedcom#NameNode.Format
[14]: https://godoc.org/github.com/elliotchance/gedcom#NameNode
[15]: https://github.com/elliotchance/gedcom/releases
