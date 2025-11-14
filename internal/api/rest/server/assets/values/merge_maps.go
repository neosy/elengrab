package avalues

import "maps"

func MergeMaps(mapsList ...map[string]string) map[string]string {
	totalLen := 0
	for _, m := range mapsList {
		totalLen += len(m)
	}

	merged := make(map[string]string, totalLen)
	for _, m := range mapsList {
		maps.Copy(merged, m)
	}

	return merged
}
