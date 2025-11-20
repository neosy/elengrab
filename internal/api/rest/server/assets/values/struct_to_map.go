package avalues

import (
	"encoding/json"
	"fmt"
)

func StructToMap(data any) map[string]any {
	var result = make(map[string]any)

	b, err := json.Marshal(data)
	if err != nil {
		return result
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(b, &temp); err != nil {
		return result
	}

	for k, v := range temp {
		result[k] = fmt.Sprintf("%v", v)
	}

	return result
}
