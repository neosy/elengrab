package avalues

import (
	"encoding/json"
	"fmt"
)

func StructToMap(data any) map[string]string {
	var result = make(map[string]string)

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
