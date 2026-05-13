package uivalues

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/neosy/elengrab/internal/pkg/httpx"
)

func mergeMaps(mapsList ...map[string]any) map[string]any {
	totalLen := 0
	for _, m := range mapsList {
		totalLen += len(m)
	}

	merged := make(map[string]any, totalLen)
	for _, m := range mapsList {
		maps.Copy(merged, m)
	}

	return merged
}

func fileNameWithHash(dir string, fileName string) (string, error) {
	filePath := filepath.Join(dir, fileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := httpx.AssetFingerprintHex32(data)
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	return fmt.Sprintf("%s.%s%s", name, hash, ext), nil
}

func fileNamesWithHash(dir string, fileNames []string) ([]string, error) {
	var result = make([]string, len(fileNames))

	for i, f := range fileNames {
		var err error
		result[i], err = fileNameWithHash(dir, f)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func fileRaw(fileName, dir string) ([]byte, error) {
	filePath := filepath.Join(dir, string(fileName))

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
