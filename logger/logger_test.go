package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ksctl/ksctl/v2/pkg/helpers/consts"
	"github.com/ksctl/ksctl/v2/pkg/types"
	"github.com/ksctl/ksctl/v2/pkg/types/controllers/cloud"
)

var (
	gL       types.LoggerFactory
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

	t.Run("Table", func(t *testing.T) {
		gL.Table(dummyCtx,
			consts.LoggingGetClusters,
			[]cloud.AllClusterData{
				{
					Name:          "fake-demo",
					CloudProvider: "fake",
					Region:        "fake-reg",
				},
			})

		gL.Table(dummyCtx, consts.LoggingGetClusters, nil)
	})

	t.Run("Box", func(t *testing.T) {
		gL.Box(dummyCtx, "Abcd", "1")
		gL.Box(dummyCtx, "Abcddedefe", "1")
		gL.Box(dummyCtx, "KUBECONFIG env var", "/jknc/csdc")
		gL.Box(dummyCtx, "KUBECONFIG env var", "jknc")
	})

	t.Run("external", func(t *testing.T) {
		gL.ExternalLogHandler(dummyCtx, consts.LogSuccess, "cdcc")
		gL.ExternalLogHandlerf(dummyCtx, consts.LogSuccess, "cdcc", "Reason", fmt.Errorf("Error"))
	})
}
