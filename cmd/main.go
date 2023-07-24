package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/danecwalker/flp-parser/pkg/defs"
	"github.com/danecwalker/flp-parser/pkg/parser"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
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
	proj, err := parser.Parse(content)
	if err != nil {
		fmt.Println("Error parsing file")
		return
	}

	for _, e := range proj.Events {
		switch e.Value().(type) {
		case string:
			if e.Kind() == defs.EventKindFilePath {
				defs.ModTextEvent(e.(*defs.TextEvent), "C:\\Users\\danew\\Downloads\\looperman-l-2217571-0289398-emeralds-lil-uzi-vert.wav")
			}
		}
	}

	flpFile2 := args[1]
	// Write file
	if err := parser.Write(proj, flpFile2); err != nil {
		fmt.Println("Error writing file")
		return
	}

	c2, err := os.ReadFile(flpFile2)
	if err != nil {
		fmt.Println("Error reading file")
		return
	}

	fmt.Println(bytes.Equal(content, c2))
}
