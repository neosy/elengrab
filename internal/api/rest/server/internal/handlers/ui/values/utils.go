package uivalues

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"
)

func MergeMaps(mapsList ...map[string]any) map[string]any {
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

func StructToMap(data any) map[string]any {
	var result = make(map[string]any)

	b, err := json.Marshal(data)
	if err != nil {
		return result
	}

	var temp map[string]any
	if err := json.Unmarshal(b, &temp); err != nil {
		return result
	}

	for k, v := range temp {
		switch vv := v.(type) {
		case bool:
			result[k] = vv
		case any:
			result[k] = vv
		default:
			result[k] = fmt.Sprintf("%v", vv)
		}
	}

	return result
}

func generateHashedFileName(dir string, fileName string) (string, error) {
	filePath := filepath.Join(dir, fileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(data))[:8] // short hash
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	return fmt.Sprintf("%s.%s%s", name, hash, ext), nil
}

func generateHashedFileNames(dir string, fileNames []string) ([]string, error) {
	var result = make([]string, len(fileNames))

	for i, f := range fileNames {
		var err error
		result[i], err = generateHashedFileName(dir, f)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
