package main

import (
	"log"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/core"
	"github.com/spf13/cobra"
)

const (
	docPath = "./docs/commands_reference"
)

func main() {
	rootCmd := core.BuildCLI()

	frontMatterFunc := func(cmd *cobra.Command) string { return "---\ntitle: \"" + cmd.CommandPath() + "\"\n---\n" }
	err := genMarkdownTree(&rootCmd, docPath, frontMatterFunc)
	if err != nil {
		log.Fatal(err)
	}
}
