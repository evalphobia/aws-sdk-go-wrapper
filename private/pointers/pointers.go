package pointers

// String returns the pointer of string.
func String(v string) *string {
	return &v
}

// Long returns converts int to int64 and returns the pointer of int64.
func Long(v int) *int64 {
	i := int64(v)
	return &i
}

// Long64 returns the pointer of int64.
func Long64(v int64) *int64 {
	return &v
}

// Float64 returns the pointer of float64.
func Float64(v float64) *float64 {
	return &v
}

// Bool returns the pointer of bool.
func Bool(b bool) *bool {
	return &b
}

// SliceString returns the slice of string pointer.
func SliceString(list []string) []*string {
	if len(list) == 0 {
		return nil
	}

	result := make([]*string, len(list))
	for i, v := range list {
		result[i] = String(v)
	}
	return result
}

// SliceFloat64 returns the slice of float64 pointer.
func SliceFloat64(list []float64) []*float64 {
	if len(list) == 0 {
		return nil
	}

	result := make([]*float64, len(list))
	for i, v := range list {
		result[i] = Float64(v)
	}
	return result
}
