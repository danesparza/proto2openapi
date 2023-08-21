package converter

import (
	"fmt"
	"github.com/emicklei/proto"
	"os"
	"strings"
)

type Converter struct {
	Source      *os.File
	PackageName string
}

// https://github.com/asaskevich/govalidator/blob/3153c74/utils.go#L101
func underscoreToCamel(in string) string {
	head := in[:1]

	repl := strings.Replace(
		strings.Title(strings.Replace(strings.ToLower(in), "_", " ", -1)),
		" ",
		"",
		-1,
	)
	return head + repl[1:]
}

func isNullable(field *proto.Field) bool {
	for _, opt := range field.Options {
		if opt.Name == "(gogoproto.nullable)" || opt.Constant.Source == "true" {
			return true
		}
	}
	return false
}

func protoTypeToOpenAPIType(protoPackage, protoType string) (typeName, typeFormat string, customType bool) {
	//	Map the proto type to the OpenAPI type
	//	Protobuf type reference: https://protobuf.dev/programming-guides/proto3/#scalar
	//	OpenAPI type reference: https://swagger.io/docs/specification/data-models/data-types/

	customType = false
	switch protoType {
	case "double":
		typeName = "number"
		typeFormat = "double"
	case "float":
		typeName = "number"
		typeFormat = "float"
	case "int32":
		typeName = "integer"
		typeFormat = "int32"
	case "int64":
		typeName = "integer"
		typeFormat = "int64"
	case "uint32":
		typeName = "number"
	case "uint64":
		typeName = "number"
	case "sint32":
		typeName = "integer"
		typeFormat = "int32"
	case "sint64":
		typeName = "integer"
		typeFormat = "int64"
	case "fixed32":
		typeName = "number"
	case "fixed64":
		typeName = "number"
	case "sfixed32":
		typeName = "integer"
		typeFormat = "int32"
	case "sfixed64":
		typeName = "integer"
		typeFormat = "int64"
	case "bool":
		typeName = "boolean"
	case "string":
		typeName = "string"
	case "bytes":
		/* Not sure the best way to represent this in OpenAPI types */
		typeName = "string"
	default:
		customType = true
		typeName = fmt.Sprintf("%s.%s", protoPackage, protoType)
	}

	return typeName, typeFormat, customType
}
