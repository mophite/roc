package fs

import (
	"errors"
	"os"
	"path"
	"strings"
)

func open(name string, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, perm)
}

func isLink(filename string) (string, error) {
	fi, err := os.Lstat(filename)
	if err != nil {
		return "", err
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		name, err := os.Readlink(filename)
		if err != nil {
			return "", err
		}
		return name, nil
	}

	return "", errors.New("not symlink")
}

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func getFilenamePrefix(s string) string {
	return strings.TrimSuffix(path.Base(s), ".log")
}
