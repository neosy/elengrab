package nfile

import "os"

func ExistFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func NotExistFile(path string) bool {
	return !ExistFile(path)
}
