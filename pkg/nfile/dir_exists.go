package nfile

import "os"

// DirExists reports whether the given path exists and is a directory.
// If the path exists but is not a directory (e.g., a file), it returns (false, nil).
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// DirNotExists reports whether the given path does not exist or is not a directory.
// If the path exists but is not a directory, it returns true.
func DirNotExists(path string) (bool, error) {
	exists, err := DirExists(path)
	return !exists, err
}
