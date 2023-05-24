package ddbxt

import (
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

// Merge accepts two map[string]types.AttributeValue and returns a new map containing the union of keys. If a key is
// present in both maps, Merge will use the value from src in the result map.
func Merge(src, dest map[string]types.AttributeValue) map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue, len(dest))
	for k, v := range dest {
		result[k] = v
	}
	for k, v := range src {
		result[k] = v
	}
	return result
}
