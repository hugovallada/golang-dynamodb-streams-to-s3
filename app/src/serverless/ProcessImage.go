package serverless

import (
	"github.com/aws/aws-lambda-go/events"
)

func ExtractVal(v events.DynamoDBAttributeValue) any {
	var val any
	switch v.DataType() {
	case events.DataTypeString:
		val = v.String()
	case events.DataTypeNumber:
		val, _ = v.Float()
	case events.DataTypeBinary:
		val = v.Binary()
	case events.DataTypeBoolean:
		val = v.Boolean()
	case events.DataTypeNull:
		val = nil
	case events.DataTypeList:
		list := []any{}
		for _, item := range v.List() {
			list = append(list, ExtractVal(item))
		}
		val = list
	case events.DataTypeMap:
		mapAttr := make(map[string]interface{}, len(v.Map()))
		for k, v := range v.Map() {
			mapAttr[k] = ExtractVal(v)
		}
		val = mapAttr
	case events.DataTypeBinarySet:
		set := [][]byte{}
		set = append(set, v.BinarySet()...)
		val = set
	case events.DataTypeNumberSet:
		set := []string{}
		set = append(set, v.NumberSet()...)
		val = set
	case events.DataTypeStringSet:
		set := []string{}
		set = append(set, v.StringSet()...)
		val = set
	}
	return val
}
