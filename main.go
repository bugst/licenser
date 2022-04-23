package main

import (
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
		fatal(3, "Ls is not a directory\n", rootDir)
	}

	license := detectLicense(rootDir)
	if onlyDetectLicense {
		fmt.Println("Detected license:")
		fmt.Println()
		for _, l := range license {
			fmt.Println(">", l)
		}
		os.Exit(0)
	}

	applyLicense(rootDir, license)
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

func applyLicense(rootDir *paths.Path, license []string) {

}
