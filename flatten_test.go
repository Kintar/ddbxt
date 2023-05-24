package ddbxt

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type flattenTestData struct {
	name     string
	path     string
	content  types.AttributeValue
	expected map[string]types.AttributeValue
}

func runFlattenTest(t *testing.T, data flattenTestData) {
	t.Helper()
	t.Run(data.name, func(t *testing.T) {
		t.Parallel()
		path := ""
		if data.name != "" {
			if data.path != "" {
				path = fmt.Sprintf("%s.%s", data.path, data.name)
			} else {
				path = data.name
			}
		}
		assert.Equal(t, data.expected, FlattenAv(data.content, path))
	})
}

type SubStruct1 struct {
	Foo string
	Bar bool
}

type SubStruct2 struct {
	Blip  []float64
	Davey byte
}

type ComplexStructForTesting struct {
	Name      string
	IntVal    int
	SubStruct SubStruct1
	ListThing []SubStruct2
}

func TestFlattenAv(t *testing.T) {
	td := ComplexStructForTesting{
		Name:   "jonathan",
		IntVal: 6512,
		SubStruct: SubStruct1{
			Foo: "absolutely",
			Bar: true,
		},
		ListThing: []SubStruct2{
			{
				Blip:  []float64{51.2, 512.94123},
				Davey: byte(15),
			},
		},
	}

	expected := make(map[string]types.AttributeValue)
	expected["Name"] = AvS(td.Name)
	expected["IntVal"] = AvN(td.IntVal)
	expected["ListThing[0].Blip[0]"] = AvN(td.ListThing[0].Blip[0])
	expected["ListThing[0].Blip[1]"] = AvN(td.ListThing[0].Blip[1])
	expected["ListThing[0].Davey"] = AvN(td.ListThing[0].Davey)
	expected["SubStruct.Foo"] = AvS(td.SubStruct.Foo)
	expected["SubStruct.Bar"] = &types.AttributeValueMemberBOOL{Value: td.SubStruct.Bar}

	m, _ := attributevalue.Marshal(td)

	assert.Equal(t, expected, FlattenAv(m, ""))
}
