package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	decode := flag.Bool("d", false, "decode a huffman-encoded file")

	inputFileName := flag.String("i", "", "name of inputFile")
	outputFileName := flag.String("o", "", "name of outputFile")

	flag.Parse()

	if !*decode && !strings.HasSuffix(*outputFileName, ".huff") {
		*outputFileName = *outputFileName + ".huff"
	}

	h := Huffman{}
	if *decode {
		ratio, err := h.Decode(*inputFileName, *outputFileName)
		if err != nil {
			fmt.Printf("error decoding file %s. Error: %v\n", *inputFileName, err)
			os.Exit(1)
		}
		fmt.Printf("Written %s. Decompression Rate: %.02f%%.\n", *outputFileName, ratio * 100)
	} else {
		ratio, err := h.Encode(*inputFileName, *outputFileName)
		if err != nil {
			fmt.Printf("error decoding file %s. Error: %v\n", *inputFileName, err)
			os.Exit(1)
		}
		fmt.Printf("Written %s. Compression Rate: %.02f%%.\n", *outputFileName, ratio * 100)
	}

	os.Exit(0)
}
