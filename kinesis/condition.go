package kinesis

// ShardIteratorType
const (
	IteratorTypeLatest      IteratorType = "LATEST"
	IteratorTypeTrimHorizon IteratorType = "TRIM_HORIZON"
)

// GetCondition has option values for `GetRecord` operation.
type GetCondition struct {
	ShardID           string
	ShardIterator     string
	ShardIteratorType IteratorType
	Limit             int64
}

// IteratorType is ShardIteratorType.
type IteratorType string

// String returns the IteratorType as string type.
// If it's empty string, then returns `LATEST`.
func (it IteratorType) String() string {
	if it.isEmpty() {
		return string(IteratorTypeLatest)
	}
	return string(it)
}

// isEmpty checks if the IteratorType is empty string or not.
func (it IteratorType) isEmpty() bool {
	return string(it) == ""
}
