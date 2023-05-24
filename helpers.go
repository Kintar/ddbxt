package ddbxt

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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

const listFmt = "%s.%s[%d]"

// FlattenAv accepts a types.AttributeValue and builds a map of its contents. The keys of the new map depend on the
// type of the incoming AttributeValue.
// - For slice types and AttributeValueMemberL, the keys are <currentPath>[index]
// - For AttributeValueMemberM, this function will recurse and the keys are <currentPath>.key[.subKey[.subKey...]]
// - For all other values, a map containing a single entry with the original value keyed as <currentPath> is returned.
func FlattenAv(v types.AttributeValue, path string) map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue)
	switch val := v.(type) {
	case *types.AttributeValueMemberM:
		for k, v := range val.Value {
			childPath := k
			if path != "" {
				childPath = fmt.Sprintf("%s.%s", path, childPath)
			}
			subResult := FlattenAv(v, childPath)
			for kk, vv := range subResult {
				result[kk] = vv
			}
		}
	case *types.AttributeValueMemberL:
		for i, v := range val.Value {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			subResult := FlattenAv(v, childPath)
			for kk, vv := range subResult {
				result[kk] = vv
			}
		}
	default:
		result[path] = v
	}

	return result
}
