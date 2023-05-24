package ddbxt

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// NumericPrimitive is a generic type constraint that is satisfied by any built-in numeric type
type NumericPrimitive interface {
	uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128 | uint | int
}

// AvN is a helper to create AttributeValueMemberN from numeric primitive types
// NOTE: This is not generally needed. attributevalue.Marshal() is what you probably want.
func AvN[T NumericPrimitive](v T) types.AttributeValue {
	return &types.AttributeValueMemberN{Value: fmt.Sprintf("%v", v)}
}

// AvNS is a helper to create AttributeValueMemberNS from numeric primitive slices
// NOTE: This is not generally needed. attributevalue.Marshal() is what you probably want.
func AvNS[T NumericPrimitive](v []T) types.AttributeValue {
	var val []string
	for _, v1 := range v {
		val = append(val, fmt.Sprintf("%v", v1))
	}
	return &types.AttributeValueMemberNS{Value: val}
}

// AvB is a helper to create AttributeValueMemberB
func AvB(v []byte) types.AttributeValue {
	return &types.AttributeValueMemberB{Value: v}
}

// AvBS is a helper to create AttributeValueMemberBS
func AvBS(v [][]byte) types.AttributeValue {
	return &types.AttributeValueMemberBS{Value: v}
}

// AvS is a helper to create AttributeValueMemberS
func AvS(v string) types.AttributeValue {
	return &types.AttributeValueMemberS{Value: v}
}

// AvL is a wrapper around attributevalue.Marshal() that swallows errors
func AvL(v []any) types.AttributeValue {
	l, _ := attributevalue.Marshal(v)
	return l
}

// AvM is a wrapper around attributevalue.Marshal() that swallows errors
func AvM(v map[string]any) types.AttributeValue {
	m, _ := attributevalue.Marshal(v)
	return m
}
