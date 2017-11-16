package rekognition

import "strings"

var msgInvalidImageEncoding = `InvalidImageFormatException: Invalid image encoding`

// IsErrorInvalidImageEncoding checks if gicen error is `Invalid image encoding`.
func IsErrorInvalidImageEncoding(err error) (ok bool) {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), msgInvalidImageEncoding)
}
