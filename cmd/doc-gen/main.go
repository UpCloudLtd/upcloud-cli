package main

import (
	"log"

	"github.com/UpCloudLtd/upcloud-cli/internal/core"

	"github.com/spf13/cobra/doc"
)

const (
	docPath = "./docs"
)

func main() {
	rootCmd := core.BuildCLI()

	err := doc.GenMarkdownTree(&rootCmd, docPath)
	if err != nil {
		log.Fatal(err)
	}
}
