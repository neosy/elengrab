package filex

import (
	"fmt"
	"os"
)

// CheckDir verifies that the given path exists and is a directory.
// Returns an error if the path does not exist or is not a directory.
func CheckDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return fmt.Errorf("cannot access directory %s: %v", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}
	return nil
}
