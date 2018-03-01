package rekognition

import "strings"

const (
	msgInvalidImageEncoding = `InvalidImageFormatException: Invalid image encoding`
	msgInvalidParameter     = `InvalidParameterException: Request has Invalid Parameters`
)

// IsErrorInvalidImageEncoding checks if gicen error is `Invalid image encoding`.
func IsErrorInvalidImageEncoding(err error) (ok bool) {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), msgInvalidImageEncoding)
}

// IsErrorInvalidParameter checks if gicen error is `InvalidParameterException`.
func IsErrorInvalidParameter(err error) (ok bool) {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), msgInvalidParameter)
}
