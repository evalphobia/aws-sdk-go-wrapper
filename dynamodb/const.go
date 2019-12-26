// DynamoDB utility

package dynamodb

import (
	SDK "github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	// attribvute types
	AttributeTypeString    = "S"
	AttributeTypeNumber    = "N"
	AttributeTypeBinary    = "B"
	AttributeTypeBool      = "BOOL"
	AttributeTypeNull      = "NULL"
	AttributeTypeMap       = "M"
	AttributeTypeList      = "L"
	AttributeTypeStringSet = "SS"
	AttributeTypeNumberSet = "NS"
	AttributeTypeBinarySet = "BS"

	conditionEQ      = "="
	conditionLE      = "<="
	conditionLT      = "<"
	conditionGE      = ">="
	conditionGT      = ">"
	conditionBETWEEN = SDK.ComparisonOperatorBetween
	conditionOR      = SDK.ConditionalOperatorOr
	conditionAND     = SDK.ConditionalOperatorAnd

	// comparison operators
	ComparisonOperatorEQ = SDK.ComparisonOperatorEq
	ComparisonOperatorNE = "NE"
	ComparisonOperatorGT = "GT"
	ComparisonOperatorLT = "LT"
	ComparisonOperatorGE = "GE"
	ComparisonOperatorLE = "LE"

	// key type name for DynamoDB Index.
	KeyTypeHash  = SDK.KeyTypeHash
	KeyTypeRange = SDK.KeyTypeRange

	SelectCount = SDK.SelectCount
)
