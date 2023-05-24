package ddbxt

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type avEqTestData struct {
	name     string
	expected bool
	v1       types.AttributeValue
	v2       types.AttributeValue
}

func makeAvEqTestData[T any](name string, v1, v2 T, equal bool) avEqTestData {
	val1, _ := attributevalue.Marshal(v1)
	val2, _ := attributevalue.Marshal(v2)
	return avEqTestData{
		name:     name,
		expected: equal,
		v1:       val1,
		v2:       val2,
	}
}

func TestAttributeValuesEqual(t *testing.T) {
	t.Parallel()

	tests := []avEqTestData{
		makeAvEqTestData("string equal", "one", "one", true),
		makeAvEqTestData("string unequal", "one", "two", false),
		makeAvEqTestData("number equal", 1, 1, true),
		makeAvEqTestData("number unequal", 2, 1, false),
		makeAvEqTestData("bool equal", true, true, true),
		makeAvEqTestData("bool unequal", false, true, false),
		makeAvEqTestData(
			"null equal",
			&types.AttributeValueMemberNULL{Value: true},
			&types.AttributeValueMemberNULL{Value: true},
			true,
		),
		makeAvEqTestData(
			"null unequal",
			&types.AttributeValueMemberNULL{Value: true},
			&types.AttributeValueMemberNULL{Value: false},
			false,
		),
		makeAvEqTestData("bytes equal", []byte{1, 2, 3, 4}, []byte{1, 2, 3, 4}, true),
		makeAvEqTestData("bytes unequal", []byte{1, 1, 1, 4}, []byte{1, 2, 3, 4}, false),
		makeAvEqTestData("string slices equal", []string{"one", "two"}, []string{"one", "two"}, true),
		makeAvEqTestData("string slices unequal", []string{"one", "two"}, []string{"one", "one"}, false),
		makeAvEqTestData("number slices equal", []float64{20.0, 25.43}, []float64{20.0, 25.43}, true),
		makeAvEqTestData("number slices unequal", []int{9, 8}, []int{8, 9}, false),
		makeAvEqTestData("bool slices equal", []bool{false, false}, []bool{false, false}, true),
		makeAvEqTestData("bool slices unequal", []bool{false, false, true}, []bool{true, false, false}, false),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, AttributeValuesEqual(test.v1, test.v2))
		})
	}
}

func TestAttributeValuesEqual_Lists(t *testing.T) {
	var s1, s2 []any
	v, _ := attributevalue.Marshal("foo")
	s1 = append(s1, v)
	s2 = append(s2, v)
	v, _ = attributevalue.Marshal(12)
	s1 = append(s1, v)
	s2 = append(s2, v)

	val1, _ := attributevalue.Marshal(s1)
	val2, _ := attributevalue.Marshal(s2)

	assert.Equal(t, true, AttributeValuesEqual(val1, val2))

	s2 = append(s2, v)
	val2, _ = attributevalue.Marshal(s2)

	assert.Equal(t, false, AttributeValuesEqual(val1, val2))
}

// TODO: I have a nagging voice that tells me this test is naive. Should probably come back to this later.
func TestAttributeValuesEqual_Maps(t *testing.T) {
	type NestedStruct struct {
		Boogie bool
	}

	type TestStruct struct {
		Name     string
		Value    int
		SubSlice []int
		Nested   NestedStruct
	}

	v1 := TestStruct{
		Name:     "v1",
		Value:    20,
		SubSlice: []int{1, 2, 3},
		Nested: NestedStruct{
			Boogie: false,
		},
	}

	m1, _ := attributevalue.Marshal(v1)
	m2, _ := attributevalue.Marshal(v1)
	assert.True(t, AttributeValuesEqual(m1, m2))

	v2 := TestStruct{
		Name:     "v1",
		Value:    20,
		SubSlice: []int{1, 3, 3},
		Nested: NestedStruct{
			Boogie: false,
		},
	}
	m2, _ = attributevalue.Marshal(v2)
	assert.False(t, AttributeValuesEqual(m1, m2))
}
