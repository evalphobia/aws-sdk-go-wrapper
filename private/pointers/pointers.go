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

// Bool returns the pointer of bool.
func Bool(b bool) *bool {
	return &b
}
