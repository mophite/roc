// Copyright (c) 2021 roc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package log

import (
    "path"
    "runtime"
    "strconv"
    "time"

    "github.com/go-roc/roc/rlog/common"
    "github.com/go-roc/roc/rlog/format"
    "github.com/go-roc/roc/rlog/output"
)

func init() {
    // Call(4) is the actual line where used
    //Overload(Call(-1))
    Overload(Call(4))
}

var defaultLogger *log

type log struct {
    opts   Option
    detail *common.Detail
}

func Overload(opts ...Options) {
    if defaultLogger != nil {
        defaultLogger = nil
    }

    defaultLogger = &log{opts: newOpts(opts...)}

    defaultLogger.detail = &common.Detail{
        Name:   defaultLogger.Name(),
        Prefix: defaultLogger.Prefix(),
    }
}

func (l *log) Fire(level, msg string) *common.Detail {
    d := *l.detail
    if l.opts.call >= 0 {
        d.Line = l.caller()
    }
    d.Timestamp = time.Now().Format(l.Formatter().Layout())
    d.Level = level
    d.Content = msg
    return &d
}

func (l *log) Name() string {
    return l.opts.name
}

func (l *log) Prefix() string {
    return l.opts.prefix
}

func (l *log) Formatter() format.Formatter {
    return l.opts.format
}

func (l *log) Output() output.Outputor {
    return l.opts.out
}

func (l *log) caller() string {
    _, file, line, ok := runtime.Caller(l.opts.call)
    //funcName := runtime.FuncForPC(pc).Name()
    if !ok {
        file = "???"
        line = 0
    }
    return path.Base(file) + ":" + strconv.Itoa(line)
}

func Close() {
    defaultLogger.opts.out.Close()
}

func Debug(content string) {
    b := defaultLogger.
        Formatter().
        Format(defaultLogger.Fire("DBUG", content))

    defaultLogger.
        Output().
        Out(common.DEBUG, b)
}

func Info(content string) {
    b := defaultLogger.
        Formatter().
        Format(defaultLogger.Fire("INFO", content))

    defaultLogger.
        Output().
        Out(common.INFO, b)
}

func Warn(content string) {
    b := defaultLogger.
        Formatter().
        Format(defaultLogger.Fire("WARN", content))

    defaultLogger.
        Output().
        Out(common.WARN, b)
}

func Error(content string) {
    b := defaultLogger.
        Formatter().
        Format(defaultLogger.Fire("ERRO", content))

    defaultLogger.
        Output().
        Out(common.ERR, b)
}

func Fatal(content string) {
    b := defaultLogger.
        Formatter().
        Format(defaultLogger.Fire("FATA", content))

    defaultLogger.
        Output().
        Out(common.FATAL, b)
}

func Stack(content string) {

    buf := make([]byte, 1<<20)
    n := runtime.Stack(buf, true)
    content += string(buf[:n]) + "\n"

    b := defaultLogger.
        Formatter().
        Format(defaultLogger.Fire("STAK", content))

    defaultLogger.
        Output().
        Out(common.STACK, b)
}
