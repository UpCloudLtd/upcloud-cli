package main

import (
	"log"

	"github.com/UpCloudLtd/cli/internal/commands/all"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/core"

	"github.com/spf13/cobra/doc"
)

func main() {
	conf := config.New()
	upctl := core.BuildRootCmd(nil, conf)
	all.BuildCommands(&upctl, conf)

	err := doc.GenMarkdownTree(&upctl, "./docs")

	if err != nil {
		log.Fatal(err)
	}
}
