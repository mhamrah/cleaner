package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var infile = flag.String("i", "", "input file to read")
	var outfile = flag.String("o", "out.dat", "output file to write")
	var l = flag.Int("l", 1024, "record length")

	flag.Parse()

	if len(*infile) == 0 || len(*outfile) == 0 || *l <= 0 {
		flag.Usage()
		os.Exit(1)
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
	count := 1

	for count > 0 {
		count, err := reader.Read(data)
		if err != nil && err != io.EOF {
			fmt.Printf("Failed to read input data stream. err: %v\n", err)
			os.Exit(1)
		}

		if count == 0 {
			writer.Flush()
			fmt.Printf("Read all bytes from %v\n", *infile)
			return
		}

		if validRecord(data, count) {
			fmt.Printf("Copying %d length records from %s to %s\n", *l, *infile, *outfile)
			out.Write(data[:count])
		} else {
			fmt.Printf("Skipping null record.\n")
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
	lookahead := count / 4
	for i := 20; i < 20+lookahead; i++ {
		if data[i] == '\x00' {
			return false
		}
	}
	return true
}
