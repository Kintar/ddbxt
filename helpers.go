package ddbxt

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// FindUpdates accepts two map[string]types.AttributeValue and returns a new map containing only those values from src
// which would update values in dest. Keys that are not shared between the two maps are ignored.
func FindUpdates(src, dest map[string]types.AttributeValue) map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue, len(src))
	// Iterate over all the members of src
	for k, v1 := range src {
		// If the value doesn't exist in dest, just skip it
		if v2, ok := dest[k]; !ok {
			continue
		} else {
			if AttributeValuesEqual(v1, v2) {
				continue
			}
		}
		result[k] = v1
	}

	return result
}

// AttributeValuesEqual compares the contents of two types.AttributeValue and returns true if they are equal.
// If both values are lists or maps with the same length, this function will recurse into their contents.
func AttributeValuesEqual(v1, v2 types.AttributeValue) bool {
	// One would think there'd be an easier way to compare the CONTENTS of a types.AttributeValue. >.<
	// Because of Amazon's design decisions here, we have to do a type switch on our src value. We then check
	// if the dest value is the same type. If it is, we compare the values. If the values match, we skip this
	// entry, otherwise we add it to the result map because either the types or the values are different.
	switch val1 := v1.(type) {
	case *types.AttributeValueMemberBOOL:
		if val2, ok := v2.(*types.AttributeValueMemberBOOL); ok {
			if val2.Value == val1.Value {
				return true
			}
		}
	case *types.AttributeValueMemberB:
		if val2, ok := v2.(*types.AttributeValueMemberB); ok {
			if bytes.Equal(val2.Value, val1.Value) {
				return true
			}
		}
	case *types.AttributeValueMemberBS:
		if val2, ok := v2.(*types.AttributeValueMemberBS); ok {
			if len(val1.Value) != len(val2.Value) {
				return false
			}
			allEquivalent := true
			for i := range val2.Value {
				allEquivalent = allEquivalent && bytes.Equal(val2.Value[i], val1.Value[i])
				if !allEquivalent {
					break
				}
			}
			return allEquivalent
		}
	case *types.AttributeValueMemberL:
		if val2, ok := v2.(*types.AttributeValueMemberL); ok {
			if len(val1.Value) != len(val2.Value) {
				return false
			}
			for i := range val1.Value {
				if !AttributeValuesEqual(val1.Value[i], val2.Value[i]) {
					return false
				}
			}
			return true
		}
	case *types.AttributeValueMemberM:
		if val2, ok := v2.(*types.AttributeValueMemberM); ok {
			if len(val1.Value) != len(val2.Value) {
				return false
			}
			allEquivalent := true
			for k, v1 := range val1.Value {
				if v2, ok := val2.Value[k]; !ok {
					return false
				} else {
					allEquivalent = allEquivalent && AttributeValuesEqual(v1, v2)
					if !allEquivalent {
						return false
					}
				}
			}
			return allEquivalent
		}
	case *types.AttributeValueMemberNULL:
		if val2, ok := v2.(*types.AttributeValueMemberNULL); ok {
			return val1.Value == val2.Value
		}
	case *types.AttributeValueMemberN:
		if val2, ok := v2.(*types.AttributeValueMemberN); ok {
			return val1.Value == val2.Value
		}
	case *types.AttributeValueMemberNS:
		if val2, ok := v2.(*types.AttributeValueMemberNS); ok {
			if len(val1.Value) != len(val2.Value) {
				return false
			}
			for i := range val2.Value {
				if val2.Value[i] != val1.Value[i] {
					return false
				}
			}
			return true
		}
	case *types.AttributeValueMemberS:
		if val2, ok := v2.(*types.AttributeValueMemberS); ok {
			return val2.Value == val2.Value
		}
	case *types.AttributeValueMemberSS:
		if val2, ok := v2.(*types.AttributeValueMemberSS); ok {
			if len(val1.Value) != len(val2.Value) {
				return false
			}
			for i := range val1.Value {
				if val1.Value[i] != val2.Value[i] {
					return false
				}
			}
		}
	}

	return false
}
