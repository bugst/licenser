//
// Copyright 2022 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/spf13/cobra"
)

var onlyDetectLicense bool

func main() {
	command := &cobra.Command{
		Long:         "licenser is a tool for keeping the license text in the source code up-to-date.",
		Use:          "licenser SOURCE_ROOT_DIR [LICENSE_FILE]",
		SilenceUsage: false,
		Example:      "licenser .",
		Run:          licenser,
	}
	command.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		fatal(1, "%s\n", err)
		return err
	})
	flags := command.Flags()
	flags.BoolVarP(&onlyDetectLicense, "detect-only", "d", false, "Only detect and print license")
	_ = command.Execute()
}

func fatal(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format, args...)
	os.Exit(code)
}

func licenser(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fatal(2, "Please specify the root directory\n")
	}
	rootDir := paths.New(args[0])
	if !rootDir.IsDir() {
		fatal(3, "%s is not a directory\n", rootDir)
	}

	var license []string
	if len(args) < 2 {
		license = detectLicense(rootDir)
	} else if l, err := paths.New(args[1]).ReadFileAsLines(); err != nil {
		fatal(4, "Error reading license file: %s\n", err)
	} else {
		license = l
	}

	if onlyDetectLicense {
		fmt.Println("Detected license:")
		fmt.Println()
		for _, l := range license {
			fmt.Println(">", l)
		}
		os.Exit(0)
	}

	sources, err := rootDir.ReadDirRecursiveFiltered(
		paths.AndFilter(
			func(file *paths.Path) bool { return file.Base() != ".git" },
		),
		paths.FilterOutDirectories(),
	)
	if err != nil {
		fatal(8, "Error reading source directory: %s\n", err)
	}

	for _, s := range sources {
		switch s.Ext() {
		case ".go", ".c", ".cpp", ".h":
			applyLicenseCStyle(s, license)
		default:
			fmt.Println("IGNORED:", s)
		}
	}
}

func detectLicense(rootDir *paths.Path) []string {
	if rootDir.Join("go.mod").Exist() {
		fmt.Println("Golang project detected")

		// First check for docs.go
		if source, err := rootDir.Join("doc.go").ReadFileAsLines(); err == nil {
			fmt.Println("Extracting license from doc.go")
			return extractLicense(source)
		}
	}

	fatal(5, "Could not find any license file in %s\n", rootDir)
	return nil
}

func extractLicense(source []string) []string {
	if len(source) == 0 {
		fatal(6, "License file is empty.\n")
	}
	for i := range source {
		source[i] = strings.TrimSpace(source[i])
	}
	if source[0] == "" {
		fatal(7, "The first line of the license file must not be empty.\n")
	}
	license := []string{}
	for len(source) > 0 && (source[0] == "//" || strings.HasPrefix(source[0], "// ")) {
		l := strings.TrimPrefix(source[0], "// ")
		l = strings.TrimPrefix(l, "//")
		license = append(license, l)
		source = source[1:]
	}
	return license
}

func applyLicenseCStyle(sourceFile *paths.Path, license []string) {
	source, err := sourceFile.ReadFileAsLines()
	if err != nil {
		fatal(8, "Error opening %s: %s\n", sourceFile, err)
	}

	output := new(bytes.Buffer)
	for _, line := range license {
		if line == "" {
			output.WriteString(fmt.Sprintln("//"))
		} else {
			output.WriteString(fmt.Sprintln("//", line))
		}
	}
	output.Write([]byte(fmt.Sprintln()))

	original := new(bytes.Buffer)
	header := true
	for _, line := range source {
		original.WriteString(fmt.Sprintln(line))
		if header && strings.HasPrefix(line, "//") {
			continue
		}
		if header {
			header = false
			if line != "" {
				// This is not a license header, we should not replace it in the output
				output.Write(original.Bytes())
			}
			continue
		}

		output.WriteString(fmt.Sprintln(line))
	}

	if bytes.Equal(original.Bytes(), output.Bytes()) {
		fmt.Println("OK", sourceFile)
		return
	}
	if err := sourceFile.WriteFile(output.Bytes()); err != nil {
		fatal(8, "Error writing %s: %s\n", sourceFile, err)
	}
	fmt.Println("UPDATED", sourceFile)
}

