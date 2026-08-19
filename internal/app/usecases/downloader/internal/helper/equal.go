package helper

func NullableValuesEqual[T comparable](v1, v2 *T) bool {
	if v1 == nil || v2 == nil {
		return v1 == v2
	}
	return *v1 == *v2
}

func ValuesEqual[T comparable](v1, v2 *T) bool {
	if v1 == nil && v2 == nil {
		return true
	}

	var vv1, vv2 T

	if v1 != nil {
		vv1 = *v1
	}
	if v2 != nil {
		vv2 = *v2
	}

	return vv1 == vv2
}
