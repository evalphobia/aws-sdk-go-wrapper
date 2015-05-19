// convert types

package sqs

func String(v string) *string {
	return &v
}

func Long(v int) *int64 {
	i := int64(v)
	return &i
}
