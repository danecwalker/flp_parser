package main

import (
	"fmt"
	"os"

	"github.com/danecwalker/flp-parser/pkg/parser"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide a flp file")
		return
	}

	flpFile := args[0]
	// Check extension of file is .flp
	if flpFile[len(flpFile)-4:] != ".flp" {
		fmt.Println("Please provide a flp file")
		return
	}

	// Check if file exists
	if _, err := os.Stat(flpFile); os.IsNotExist(err) {
		fmt.Println("File does not exist")
		return
	}

	// Read file
	content, err := os.ReadFile(flpFile)
	if err != nil {
		fmt.Println("Error reading file")
		return
	}

	// Parse file
	_, err = parser.Parse(content)
	if err != nil {
		fmt.Println("Error parsing file")
		return
	}
}
