// gen-docs generates man pages for the GrooveKit CLI using cobra/doc.
// Usage: go run ./tools/gen-docs [output-dir] [version]
package main

import (
	"log"
	"os"

	"github.com/scookdev/groovekit-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	outDir := "./man1"
	if len(os.Args) > 1 {
		outDir = os.Args[1]
	}

	version := "dev"
	if len(os.Args) > 2 {
		version = os.Args[2]
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	root := cmd.RootCmd()
	disableAutoGenTag(root)

	header := &doc.GenManHeader{
		Title:   "GROOVEKIT",
		Section: "1",
		Source:  "GrooveKit " + version,
		Manual:  "GrooveKit CLI Manual",
	}

	if err := doc.GenManTree(root, header, outDir); err != nil {
		log.Fatalf("failed to generate man pages: %v", err)
	}
}

// disableAutoGenTag removes the auto-generated timestamp tag from all commands
// so that man pages are reproducible across builds.
func disableAutoGenTag(c *cobra.Command) {
	c.DisableAutoGenTag = true
	for _, sub := range c.Commands() {
		disableAutoGenTag(sub)
	}
}
