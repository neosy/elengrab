package hash

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/assetx"
)

func isPointer(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Ptr
}

func deref(v any) any {
	rv := reflect.ValueOf(v)

	for rv.IsValid() && rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	return rv.Interface()
}

func MetaHashHex32(values []any) string {
	var hashValues []string

	for _, v := range values {
		var value string
		switch v := v.(type) {
		case string:
			value = v
		case *string:
			if v == nil {
				continue
			}
			value = *v
		case time.Time:
			value = v.String()
		case *time.Time:
			if v == nil {
				continue
			}
			value = v.String()
		case uuid.UUID:
			value = v.String()
		case *uuid.UUID:
			if v == nil {
				continue
			}
			value = v.String()
		default:
			if isPointer(v) {
				if v != nil {
					vv := deref(v)
					if vv != nil {
						value = fmt.Sprint(vv)
					}
				}
			} else {
				value = fmt.Sprint(v)
			}
		}
		hashValues = append(hashValues, value)
	}

	hashDataRaw := []byte(strings.Join(hashValues, ""))
	return assetx.AssetFingerprintHex32(hashDataRaw)
}
