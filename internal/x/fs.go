package x

import (
	"os"
	"path/filepath"
	"strings"
)

const pathSeparator = string(os.PathSeparator)

func GetLastPwd() string {
	f := func(s string, pos, length int) string {
		runes := []rune(s)
		l := pos + length
		if l > len(runes) {
			l = len(runes)
		}
		return string(runes[pos:l])
	}
	directory := GetPwd()
	return f(directory, 0, strings.LastIndex(directory, pathSeparator))
}

func GetPwd() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func GetProjectName() string {
	f := func(s string, pos int) string {
		runes := []rune(s)
		return string(runes[pos:])
	}

	directory := GetPwd()
	return strings.Trim(f(directory, strings.LastIndex(directory, pathSeparator)), pathSeparator)
}
