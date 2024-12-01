package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ksctl/cli/cli/cmd"
	"github.com/spf13/cobra/doc"
)

func filePrepender(filename string) string {
	cmdName := filepath.Base(filename)
	cmdName = cmdName[:len(cmdName)-3] // Remove .md extension

	return fmt.Sprintf(`---
title: %s
description: Command documentation for %s
---

`, cmdName, cmdName)
}

func linkHandler(name string) string {
	return name
}

func main() {
	outputDir := "./gen/docs.md"

	if err := doc.GenMarkdownTreeCustom(cmd.RootCmd, outputDir, filePrepender, linkHandler); err != nil {
		log.Fatal(err)
	}
}
