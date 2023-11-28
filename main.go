package main

import (
	"fmt"
	"huffman-coding/pkg/huffman"
	"os"
)

func PrintUsage() {
	fmt.Println("usage: ./main [-d] <input_path> <output_path>")
	fmt.Println("use option -d to decode input file")
}
func main() {

	args := os.Args[1:]
	
	if len(args) == 0 {
		fmt.Println("not enough arguments.")
		PrintUsage()
		os.Exit(1)
	}

	decode := false
	if args[0][0] == '-' {
		if args[0] == "-d" {
			decode = true
		} else {
			fmt.Println("unknown command: ", args[0])
			PrintUsage()
			os.Exit(1)
		}
	}

	if decode {
		if len(args) < 3 {
			fmt.Println("not enough arguments.")
			PrintUsage()
			os.Exit(1)
		}
	} else {
		if len(args) < 2 {
			fmt.Println("not enough arguments.")
			PrintUsage()
			os.Exit(1)
		}
	}
	h := huffman.Huffman{
		Debug: false,
	}

	if decode {
		inputFile := args[1]
		outputFile := args[2]
		_, err := h.Decode(inputFile, outputFile)
		if(err != nil) {
			fmt.Print("decode error ", err)
			os.Exit(1)
		}
	}else{
		inputFile := args[0]
		outputFile := args[1]
		_, err := h.Encode(inputFile, outputFile)
		if(err != nil) {
			fmt.Print("encode error ", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
