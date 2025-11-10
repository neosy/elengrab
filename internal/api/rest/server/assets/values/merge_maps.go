package avalues

func MergeMaps(maps ...map[string]string) map[string]string {
	totalLen := 0
	for _, m := range maps {
		totalLen += len(m)
	}

	merged := make(map[string]string, totalLen)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}

	return merged
}
