package main

import (
	"log"

	"github.com/UpCloudLtd/cli/internal/commands/all"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/core"

	"github.com/spf13/cobra/doc"
)

const (
	docPath = "./docs"
)

func main() {
	conf := config.New()
	upctl := core.BuildRootCmd(nil, conf)
	all.BuildCommands(&upctl, conf)

	err := doc.GenMarkdownTree(&upctl, docPath)

	if err != nil {
		log.Fatal(err)
	}
}
