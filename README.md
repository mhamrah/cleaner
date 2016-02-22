# Cleaner

A go program to copy only valid byte-level records from one file to another.

## Install

If you have [go](golang.org) installed and configured, simply:

```
go get github.com/mhamrah/cleaner
```

You can compile binaries for other platforms via:

```
GOARCH=amd64 GOOS=freebsd go build main.go
```

## Usage

Flags:

* -i, the input file
* -o, the outfile file (default out.dat)
* -l, the total record length, including headers (default 1024)


