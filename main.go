package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	decode := flag.Bool("d", false, "specifies that program should decode data")
	input := flag.String("input", "", "input file")
	output := flag.String("output", "", "output file")
	flag.Parse()

	inputFile, err := os.Open(*input)
	if err != nil {
		fmt.Println(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(*output)
	if err != nil {
		fmt.Println(err)
	}
	defer outputFile.Close()

	if !*decode {
		w := NewWriter(outputFile)
		data, err := io.ReadAll(inputFile)
		if err != nil {
			fmt.Println(err)
		}
		if _, err = w.Write(data); err != nil {
			fmt.Println(err)
		}
		if err = w.Close(); err != nil {
			fmt.Println(err)
		}
	} else {
		r := NewReader(inputFile)
		data, err := io.ReadAll(r)
		if _, err = outputFile.Write(data); err != nil {
			fmt.Println(err)
		}
	}
}
