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
	"encoding/json"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"sync"

	box "github.com/Delta456/box-cli-maker/v2"
	"github.com/fatih/color"
	"github.com/ksctl/ksctl/v2/pkg/logger"
	"github.com/rodaine/table"

	"time"
)

type GeneralLog struct {
	mu      *sync.Mutex
	writter io.Writer
	level   uint
}

func (l *GeneralLog) ExternalLogHandler(ctx context.Context, msgType logger.CustomExternalLogLevel, message string) {
	l.log(false, false, ctx, msgType, message)
}

func (l *GeneralLog) ExternalLogHandlerf(ctx context.Context, msgType logger.CustomExternalLogLevel, format string, args ...interface{}) {
	l.log(false, false, ctx, msgType, format, args...)
}

func formGroups(disableContext bool, ctx context.Context, v ...any) (format string, vals []any) {
	if len(v) == 0 {
		return "\n", nil
	}
	_format := strings.Builder{}

	defer func() {
		format = _format.String()
	}()
	if !disableContext {
		_format.WriteString(color.HiBlackString("component=") + "%s ")
		vals = append(vals, getPackageName(ctx))
	}
	i := 0
	for ; i+1 < len(v); i += 2 {
		if !reflect.TypeOf(v[i+1]).Implements(reflect.TypeOf((*error)(nil)).Elem()) &&
			(reflect.TypeOf(v[i+1]).Kind() == reflect.Interface ||
				reflect.TypeOf(v[i+1]).Kind() == reflect.Ptr ||
				reflect.TypeOf(v[i+1]).Kind() == reflect.Struct) {
			_format.WriteString(fmt.Sprintf("%s", v[i]) + "=%#v ")
		} else {
			_format.WriteString(fmt.Sprintf("%s", v[i]) + "=%v ")
		}

		vals = append(vals, v[i+1])
	}

	for ; i < len(v); i++ {
		_format.WriteString("!!EXTRA:%v ")
		vals = append(vals, v[i])
	}
	_format.WriteString("\n")
	return
}

func isLogEnabled(level uint, msgType logger.CustomExternalLogLevel) bool {
	if msgType == LogDebug {
		return level >= 9
	}
	return true
}

func (l *GeneralLog) logErrorf(disableContext bool, disablePrefix bool, ctx context.Context, msg string, args ...any) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !disablePrefix {
		prefix := fmt.Sprintf("%s%s ", getTime(l.level), LogError)
		msg = prefix + msg
	}
	format, _args := formGroups(disableContext, ctx, args...)

	var errMsg error
	if _args == nil {
		errMsg = fmt.Errorf(msg + " " + format)
	} else {
		errMsg = fmt.Errorf(msg+" "+format, _args...)
	}

	return errMsg
}

func (l *GeneralLog) log(disableContext bool, useGroupFormer bool, ctx context.Context, msgType logger.CustomExternalLogLevel, msg string, args ...any) {
	if !isLogEnabled(l.level, msgType) {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	prefix := fmt.Sprintf("%s%s ", getTime(l.level), msgType)

	if useGroupFormer {
		msgColored := ""
		switch msgType {
		case LogSuccess:
			msgColored = color.HiGreenString(msg)
		case LogWarning:
			msgColored = color.HiYellowString(msg)
		case LogDebug:
			msgColored = color.HiMagentaString(msg)
		case LogNote:
			msgColored = color.HiCyanString(msg)
		case LogInfo:
			msgColored = color.HiBlueString(msg)
		}
		msg = prefix + msgColored
		format, _args := formGroups(disableContext, ctx, args...)
		if _args == nil {
			if disableContext && msgType == LogError {
				l.boxBox(
					"Error", msg+" "+format, "Red")
				return
			}
			fmt.Fprint(l.writter, msg+" "+format)
		} else {
			if disableContext && msgType == LogError {
				l.boxBox(
					"Error", fmt.Sprintf(msg+" "+format, _args...), "Red")
				return
			}
			fmt.Fprintf(l.writter, msg+" "+format, _args...)
		}
	} else {
		args = append([]any{getPackageName(ctx)}, args...)
		fmt.Fprintf(l.writter, prefix+color.HiBlackString("component=")+"%s "+msg+"\n", args...)
	}
}

func getTime(level uint) string {
	t := time.Now()
	return color.HiBlackString(fmt.Sprintf("%02d:%02d:%02d ", t.Hour(), t.Minute(), t.Second()))
}

func NewLogger(verbose int, out io.Writer) *GeneralLog {

	var ve uint

	if verbose < 0 {
		ve = 9
	}

	return &GeneralLog{
		writter: out,
		level:   ve,
		mu:      new(sync.Mutex),
	}
}

var (
	LogWarning = logger.CustomExternalLogLevel(color.New(color.FgBlack, color.BgYellow).Sprintf("WARN"))
	LogInfo    = logger.CustomExternalLogLevel(color.New(color.FgBlack, color.BgBlue).Sprintf("INFO"))
	LogNote    = logger.CustomExternalLogLevel(color.New(color.FgBlack, color.BgCyan).Sprintf("NOTE"))
	LogDebug   = logger.CustomExternalLogLevel(color.New(color.FgBlack, color.BgMagenta).Sprintf("DEBUG"))
	LogSuccess = logger.CustomExternalLogLevel(color.New(color.FgBlack, color.BgGreen).Sprintf("SUCC"))
	LogError   = logger.CustomExternalLogLevel(color.New(color.FgBlack, color.BgRed).Sprintf("ERR"))
)

func (l *GeneralLog) Print(ctx context.Context, msg string, args ...any) {
	l.log(false, true, ctx, LogInfo, msg, args...)
}

func (l *GeneralLog) Success(ctx context.Context, msg string, args ...any) {
	l.log(false, true, ctx, LogSuccess, msg, args...)
}

func (l *GeneralLog) Note(ctx context.Context, msg string, args ...any) {
	l.log(false, true, ctx, LogNote, msg, args...)
}

func (l *GeneralLog) Debug(ctx context.Context, msg string, args ...any) {
	l.log(false, true, ctx, LogDebug, msg, args...)
}

func (l *GeneralLog) Error(msg string, args ...any) {
	l.log(true, true, nil, LogError, msg, args...)
}

func (l *GeneralLog) NewError(ctx context.Context, msg string, args ...any) error {
	return l.logErrorf(false, true, ctx, msg, args...)
}

func (l *GeneralLog) Warn(ctx context.Context, msg string, args ...any) {
	l.log(false, true, ctx, LogWarning, msg, args...)
}

func (l *GeneralLog) Table(ctx context.Context, op logger.LogClusterDetail, data []logger.ClusterDataForLogging) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	if op == logger.LoggingGetClusters {
		tbl := table.New("ClusterName", "Region", "ClusterType", "CloudProvider", "BootStrap", "WorkerPlaneNodes", "ControlPlaneNodes", "EtcdNodes", "CloudManagedNodes")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, row := range data {
			tbl.AddRow(row.Name,
				row.Region,
				string(row.ClusterType),
				string(row.CloudProvider),
				string(row.K8sDistro),
				row.NoWP, row.NoCP, row.NoDS, row.NoMgt,
			)
		}

		println()
		tbl.Print()
		println()
		println()
	} else if op == logger.LoggingInfoCluster {
		a, err := json.MarshalIndent(data[0], "", " ")
		if err != nil {
			panic(err)
		}
		l.Box(ctx, "Cluster Data", string(a))
	}

}

func (l *GeneralLog) boxBox(title, lines string, color string) {

	px := 4

	if len(title) >= 2*px+len(lines) {
		// some maths
		px = int(math.Ceil(float64(len(title)-len(lines))/2)) + 1
	}

	Box := box.New(box.Config{
		Px:       px,
		Py:       2,
		Type:     "Round",
		TitlePos: "Top",
		Color:    color})

	Box.Println(title, addLineTerminationForLongStrings(lines))
}

func (l *GeneralLog) Box(ctx context.Context, title string, lines string) {

	l.Debug(ctx, "PostUpdate Box", "title", len(title), "lines", len(lines))

	l.boxBox(title, lines, "Yellow")
}
