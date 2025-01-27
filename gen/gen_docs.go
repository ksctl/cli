// Copyright 2025 Ksctl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// import (
// 	"fmt"
// 	"log"
// 	"path/filepath"

// 	"github.com/ksctl/cli/cli/cmd"
// 	"github.com/spf13/cobra/doc"
// )

// func filePrepender(filename string) string {
// 	cmdName := filepath.Base(filename)
// 	cmdName = cmdName[:len(cmdName)-3] // Remove .md extension

// 	return fmt.Sprintf(`---
// title: %s
// description: Command documentation for %s
// ---

// `, cmdName, cmdName)
// }

// func linkHandler(name string) string {
// 	return name
// }

// func main() {
// 	outputDir := "./gen/docs.md"

// 	if err := doc.GenMarkdownTreeCustom(cmd.RootCmd, outputDir, filePrepender, linkHandler); err != nil {
// 		log.Fatal(err)
// 	}
// }
