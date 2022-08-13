<h1 align=center><code>c4udit</code></h1>

## Introduction

`c4udit` is a static analyzer for solidity contracts based on regular
expressions specifically crafted for [Code4Rena](https://code4rena.com) contests.

It is capable of finding low risk issues and gas optimization documented in the
[c4-common-issues](https://github.com/byterocket/c4-common-issues) repository.

Note that `c4udit` uses [c4-common-issues](https://github.com/byterocket/c4-common-issues)'s issue identifiers.

## Installation

First you need to have the Go toolchain installed. You can find instruction [here](https://go.dev/doc/install).

Then install `c4udit` with:
```
$ go install github.com/byterocket/c4udit@latest
```

To just build the binary:
```
$ git clone https://github.com/Ellahinator/c4udit.git
$ cd c4udit/
$ go build .
```
Now you should be able to run `c4udit` with:
```
$ ./c4udit
```

## Usage

```
Usage:
	c4udit [flags] [files...]

Flags:
	-h    Print help text.
	-s    Save report as file.
	-t    Add ToC to file.
```

## Example

Running `c4udit` against dummy.sol:
```
$ ./c4udit -s dummy.sol
$ ./c4udit -t
