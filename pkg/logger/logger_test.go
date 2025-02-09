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

package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ksctl/ksctl/v2/pkg/logger"
)

var (
	gL       logger.Logger
	dummyCtx = context.TODO()
)

func TestMain(m *testing.M) {
	gL = NewLogger(-1, os.Stdout)
	_ = NewLogger(0, os.Stdout)
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestHelperToAddLineTerminationForLongStrings(t *testing.T) {
	test := fmt.Sprintf("Argo Rollouts (Ver: %s) is a Kubernetes controller and set of CRDs which provide advanced deployment capabilities such as blue-green, canary, canary analysis, experimentation, and progressive delivery features to Kubernetes.", "v1.2.4")

	x := strings.Split(addLineTerminationForLongStrings(test), "\n")
	for _, line := range x {
		if len(line) > limitCol+1 {
			t.Errorf("Line too long: %s, got: %d, expected: %d", line, len(line), limitCol)
		}
	}
}

func TestPrinters(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		gL.Success(dummyCtx, "FAKE", "type", "success")
	})

	t.Run("Warn", func(t *testing.T) {
		gL.Warn(dummyCtx, "FAKE", "type", "warn")
	})

	t.Run("Error", func(t *testing.T) {
		gL.Error("FAKE", "type", "error")
	})

	t.Run("Debug", func(t *testing.T) {
		gL.Debug(dummyCtx, "FAKE", "type", "debugging")
	})

	t.Run("Note", func(t *testing.T) {
		gL.Note(dummyCtx, "FAKE", "type", "note")
	})

	t.Run("Print", func(t *testing.T) {
		gL.Print(dummyCtx, "FAKE", "type", "print")
	})

	t.Run("Box", func(t *testing.T) {
		gL.Box(dummyCtx, "Abcd", "1")
		gL.Box(dummyCtx, "Abcddedefe", "1")
		gL.Box(dummyCtx, "KUBECONFIG env var", "/jknc/csdc")
		gL.Box(dummyCtx, "KUBECONFIG env var", "jknc")
	})

	t.Run("external", func(t *testing.T) {
		gL.ExternalLogHandler(dummyCtx, logger.LogSuccess, "cdcc")
		gL.ExternalLogHandlerf(dummyCtx, logger.LogSuccess, "cdcc", "Reason", fmt.Errorf("Error"))
	})
}
