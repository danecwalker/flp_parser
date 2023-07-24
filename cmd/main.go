package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	if proj.Events[0].Kind() != defs.EventKindFLP_Version {
		fmt.Println("Error parsing file")
		return
	}

	fmt.Println("FLP Version:", proj.Events[0].Value().(string))

	// parser.ResolveFactory(proj)

	for _, ev := range proj.Events[1:] {
		switch ev.Kind() {
		case defs.EventKindFilePath:
			fmt.Println("File Path:", ev.Value().(string))
			os.Mkdir("output", os.ModePerm)
			f, err := os.Create("output/" + filepath.Base(strings.ReplaceAll(ev.Value().(string), "\x00", "")))
			if err != nil {
				fmt.Println("Error creating file", err)
				continue
			}

			src, err := os.ReadFile(ev.Value().(string))
			if err != nil {
				fmt.Println("Error reading file")
				continue
			}

			f.Write(src)
			f.Close()
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
