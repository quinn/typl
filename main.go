package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/quinn/typl/gen"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a glob pattern as an argument")
		return
	}

	globPattern := os.Args[1]
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		fmt.Printf("Error matching glob pattern: %v\n", err)
		return
	}

	if len(matches) == 0 {
		fmt.Println("No files found matching the glob pattern")
		return
	}

	for _, templatePath := range matches {
		packageName := filepath.Base(filepath.Dir(templatePath))
		outputPath := strings.TrimSuffix(templatePath, filepath.Ext(templatePath)) + ".go"
		err := gen.Exec(templatePath, outputPath, packageName)
		if err != nil {
			fmt.Printf("Error generating structs for %s: %v\n", templatePath, err)
		} else {
			fmt.Printf("Generated %s\n", outputPath)
		}
	}
}
