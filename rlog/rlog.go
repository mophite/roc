package rlog

import (
	"fmt"

	"roc/rlog/log"
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
