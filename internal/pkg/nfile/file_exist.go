package nfile

import "os"

// FileExists reports whether the given path exists and is a regular file.
// If the path exists but is a directory, it returns (false, nil).
func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.Mode().IsRegular(), nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// FileNotExists reports whether the given path does not exist or is not a regular file.
// If the path exists but is not a regular file (e.g., a directory or symlink), it returns true.
func FileNotExists(path string) (bool, error) {
	exists, err := FileExists(path)
	return !exists, err
}
