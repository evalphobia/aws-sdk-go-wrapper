package pointers

func String(v string) *string {
	return &v
}

func Long(v int) *int64 {
	i := int64(v)
	return &i
}

func Bool(b bool) *bool {
	return &b
}
