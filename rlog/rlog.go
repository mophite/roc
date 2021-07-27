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

package rlog

import (
	"fmt"

	"github.com/go-roc/roc/rlog/log"
)

func Debug(msg ...interface{}) {

	log.Debug(fmt.Sprintln(msg...))
}

func Info(msg ...interface{}) {
	log.Info(fmt.Sprintln(msg...))
}

func Warn(msg ...interface{}) {
	log.Warn(fmt.Sprintln(msg...))
}

func Error(msg ...interface{}) {
	log.Error(fmt.Sprintln(msg...))
}

func Fatal(msg ...interface{}) {
	log.Fatal(fmt.Sprintln(msg...))
}

func Stack(msg ...interface{}) {
	log.Stack(fmt.Sprintln(msg...))
}

func Debugf(f string, msg ...interface{}) {

	log.Debug(fmt.Sprintf(f+"\n", msg...))
}

func Infof(f string, msg ...interface{}) {
	log.Info(fmt.Sprintf(f+"\n", msg...))
}

func Warnf(f string, msg ...interface{}) {
	log.Warn(fmt.Sprintf(f+"\n", msg...))
}

func Errorf(f string, msg ...interface{}) {
	log.Error(fmt.Sprintf(f+"\n", msg...))
}

func Fatalf(f string, msg ...interface{}) {
	log.Fatal(fmt.Sprintf(f+"\n", msg...))
}

func Stackf(f string, msg ...interface{}) {
	log.Stack(fmt.Sprintf(f+"\n", msg...))
}
