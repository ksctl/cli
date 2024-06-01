package main

import (
	"log"

	"github.com/ksctl/cli/cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {

	if err := doc.GenReSTTree(cmd.RootCmd, "./gen/docs.rst"); err != nil {
		log.Fatal(err)
	}

	if err := doc.GenMarkdownTree(cmd.RootCmd, "./gen/docs.md"); err != nil {
		log.Fatal(err)
	}
}
