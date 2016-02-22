package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var c *int

func main() {
	var infile = flag.String("i", "", "input file to read")
	var outfile = flag.String("o", "out.dat", "output file to write")
	var l = flag.Int("l", 1024, "record length including header, must be > 20")
	c = flag.Int("c", 0, "lookahead length excluding header must be <= record length - 20")

	flag.Parse()

	if len(*infile) == 0 || len(*outfile) == 0 || *l <= 20 || *c > *l-20 {
		flag.Usage()
		os.Exit(1)
	}

	if *c == 0 {
		*c = *l - 20
	}

	in, err := os.Open(*infile)
	if err != nil {
		fmt.Printf("Can't read %s, err: %v", *infile, err)
		os.Exit(1)
	}

	out, err := os.Create(*outfile)
	if err != nil {
		fmt.Printf("Can't open %s for writing, err: %v\n", *infile, err)
		os.Exit(1)
	}

	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	data := make([]byte, *l)
	nullCount, copyCount := 0, 0

	byteCount, err := reader.Read(data)
	if err != nil && err != io.EOF {
		fmt.Printf("Failed to copy control record. err:", err)
		os.Exit(1)
	}

	out.Write(data[:byteCount])

	for byteCount > 0 {
		byteCount, err := reader.Read(data)
		if err != nil && err != io.EOF {
			fmt.Printf("Failed to read input data stream. err: %v\n", err)
			os.Exit(1)
		}

		if byteCount == 0 {
			writer.Flush()
			fmt.Printf("Read all bytes from file %v. Encounterd %v null records and copied %v.\n", *infile, nullCount, copyCount)
			return
		}

		if validRecord(data, byteCount) {
			out.Write(data[:byteCount])
			copyCount += 1
		} else {
			nullCount += 1
		}

		if (nullCount+copyCount)%1000 == 0 {
			fmt.Printf("In progress. Copied %v records and encountered %v null records.\n", copyCount, nullCount)
		}
	}

}

func validRecord(data []byte, count int) bool {
	if len(data) < 20 || count < 20 {
		return false
	}
	if data[0] == '\x00' {
		return false
	}
	for i := 20; i < 20+*c; i++ {
		if data[i] != '\x20' {
			return true
		}
	}

	return false
}
