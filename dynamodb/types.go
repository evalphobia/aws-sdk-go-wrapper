// convert types

package dynamodb

func String(v string) *string {
	return &v
}

func Boolean(v bool) *bool {
	return &v
}

func Long(v int64) *int64 {
	return &v
}
