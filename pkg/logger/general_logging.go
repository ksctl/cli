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
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
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
	l.log(false, msgType, message)
}

func (l *GeneralLog) ExternalLogHandlerf(ctx context.Context, msgType logger.CustomExternalLogLevel, format string, args ...interface{}) {
	l.log(false, msgType, format, args...)
}

func formGroups(v ...any) (format string, vals []any) {
	if len(v) == 0 {
		return "\n", nil
	}
	_format := strings.Builder{}

	defer func() {
		format = _format.String()
	}()

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
	if msgType == logger.LogDebug {
		return level >= 9
	}
	return true
}

func (l *GeneralLog) logErrorf(disablePrefix bool, msg string, args ...any) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !disablePrefix {
		prefix := fmt.Sprintf("%s%s ", getTime(), logger.LogError)
		msg = prefix + msg
	}
	format, _args := formGroups(args...)

	var errMsg error
	if _args == nil {
		errMsg = fmt.Errorf("%s %s", msg, format)
	} else {
		errMsg = fmt.Errorf(msg+" "+format, _args...)
	}

	return errMsg
}

func (l *GeneralLog) log(useGroupFormer bool, msgType logger.CustomExternalLogLevel, msg string, args ...any) {
	if !isLogEnabled(l.level, msgType) {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	prefix := fmt.Sprintf("%s%s ", getTime(), msgType)

	if useGroupFormer {

		msgColored := ""
		switch msgType {
		case logger.LogSuccess:
			msgColored = color.HiGreenString(msg)
		case logger.LogWarning:
			msgColored = color.HiYellowString(msg)
		case logger.LogDebug:
			msgColored = color.HiMagentaString(msg)
		case logger.LogNote:
			msgColored = color.HiCyanString(msg)
		case logger.LogInfo:
			msgColored = color.HiBlueString(msg)
		case logger.LogError:
			msgColored = color.HiRedString(msg)
		}
		msg = prefix + msgColored
		format, _args := formGroups(args...)
		if _args == nil {
			if msgType == logger.LogError {
				l.boxBox(
					"🛑 We Have Problem", msgColored+" "+format, "Red")
				return
			}
			fmt.Fprint(l.writter, msg+" "+format)
		} else {
			if msgType == logger.LogError {
				l.boxBox(
					"🛑 We Have Problem", fmt.Sprintf(msgColored+" "+format, _args...), "Red")
				return
			}
			fmt.Fprintf(l.writter, msg+" "+format, _args...)
		}
	} else {
		fmt.Fprintf(l.writter, prefix+msg+"\n", args...)
	}
}

func getTime() string {
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

func (l *GeneralLog) Print(ctx context.Context, msg string, args ...any) {
	l.log(true, logger.LogInfo, msg, args...)
}

func (l *GeneralLog) Success(ctx context.Context, msg string, args ...any) {
	l.log(true, logger.LogSuccess, msg, args...)
}

func (l *GeneralLog) Note(ctx context.Context, msg string, args ...any) {
	l.log(true, logger.LogNote, msg, args...)
}

func (l *GeneralLog) Debug(ctx context.Context, msg string, args ...any) {
	l.log(true, logger.LogDebug, msg, args...)
}

func (l *GeneralLog) Error(msg string, args ...any) {
	l.log(true, logger.LogError, msg, args...)
}

func (l *GeneralLog) NewError(ctx context.Context, msg string, args ...any) error {
	return l.logErrorf(true, msg, args...)
}

func (l *GeneralLog) Warn(ctx context.Context, msg string, args ...any) {
	l.log(true, logger.LogWarning, msg, args...)
}

func (l *GeneralLog) Table(ctx context.Context, headers []string, data [][]string) {
	headerFmt := color.New(color.FgHiBlack, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgHiGreen).SprintfFunc()

	var dataToPrint [][]interface{} = make([][]interface{}, 0, len(data))
	for _, v := range data {
		var row []interface{}
		for _, vv := range v {
			row = append(row, vv)
		}
		dataToPrint = append(dataToPrint, row)
	}

	var header []interface{}
	for _, v := range headers {
		header = append(header, v)
	}

	tbl := table.New(header...)
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, row := range dataToPrint {
		tbl.AddRow(row...)
	}
	tbl.Print()
}

func (l *GeneralLog) boxBox(title, lines string, colorName string) {
	var borderColor lipgloss.Color
	switch colorName {
	case "Red":
		borderColor = lipgloss.Color("9") // Bright red
	case "Green":
		borderColor = lipgloss.Color("10") // Bright green
	default:
		borderColor = lipgloss.Color("#555555") // Default gray
	}

	width := max(min(len(lines), 80), min(len(title), 80)) + 2

	var builder strings.Builder
	builder.WriteString("\n\n")

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Align(lipgloss.Center).
		Padding(0, 1).
		Width(width)

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Align(lipgloss.Center).
		PaddingLeft(1).
		PaddingRight(1)

	titleStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Align(lipgloss.Center).
		Bold(true).
		PaddingBottom(1).
		PaddingLeft(1).
		PaddingRight(1)

	cardContent := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(title),
		contentStyle.Render(lines),
	)

	builder.WriteString(boxStyle.Render(cardContent))
	builder.WriteString("\n\n")

	fmt.Fprintln(l.writter, builder.String())
}

func (l *GeneralLog) Box(ctx context.Context, title string, lines string) {

	l.Debug(ctx, "PostUpdate Box", "title", len(title), "lines", len(lines))

	l.boxBox(title, lines, "Green")
}
