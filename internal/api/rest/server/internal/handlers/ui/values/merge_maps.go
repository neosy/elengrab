package uivalues

import "maps"

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
