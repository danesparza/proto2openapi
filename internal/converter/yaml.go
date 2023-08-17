package converter

import (
	"fmt"
	"github.com/emicklei/proto"
	"github.com/rs/zerolog/log"
)

// ConvertToYAML converts the Source proto file to YAML format
// and returns it as a string
func (c *Converter) ConvertToYAML() (string, error) {
	retval := ""

	// Parse
	parser := proto.NewParser(c.Source)
	definition, err := parser.Parse()
	if err != nil {
		log.Err(err).Str("proto", c.Source.Name()).Msg("Problem parsing the proto file")
		return retval, err
	}

	//	Get package information
	var pkg *proto.Package
	for _, elem := range definition.Elements {

		//	Get message information
		pkgInfo, ok := elem.(*proto.Package)
		if !ok {
			continue
		}

		pkg = pkgInfo
	}
	c.PackageName = pkg.Name

	//	Indicate that the file was parsed successfully
	log.Info().Str("proto", c.Source.Name()).Str("package", c.PackageName).Msg("Parsed proto file")

	//	Get the preamble:
	retval += c.GetYAMLPreamble()

	//	Start parsing the messages:
	for _, elem := range definition.Elements {

		//	Get message information
		message, ok := elem.(*proto.Message)
		if !ok {
			continue
		}

		retval += c.GetYAMLTypeFromMessage(message)
	}

	return retval, nil
}

// GetYAMLPreamble emits the OpenAPI YAML schema preamble
func (c *Converter) GetYAMLPreamble() string {
	return fmt.Sprint("\tschemas:\n")
}

// GetYAMLTypeFromMessage emits OpenAPI YAML type information for a proto message
func (c *Converter) GetYAMLTypeFromMessage(message *proto.Message) string {
	retval := ""

	//	Write the basic type information
	retval += fmt.Sprintf("\t\t%s:\n\t\t\ttype: object\n", message.Name)

	//	If this has fields, write 'properties' and then write the fields under that.
	if len(message.Elements) > 0 {
		retval += fmt.Sprintf("\t\t\tproperties:\n")
	}

	// Write out each of the fields
	for _, node := range message.Elements {

		//	Get field information
		field, ok := node.(*proto.NormalField)
		if !ok {
			continue
		}

		//	Format and emit the name
		name := underscoreToCamel(field.Name)
		retval += fmt.Sprintf("\t\t\t\t%s:", name)

		//	Map the proto type to the OpenAPI type
		//	OpenAPI type reference: https://swagger.io/docs/specification/data-models/data-types/
		//	Protobuf type reference: https://protobuf.dev/programming-guides/proto3/#scalar
		var t string
		var tf string
		switch field.Type {
		case "bytes":
			fallthrough
		case "string":
			t = "string"
		case "uint64":
			t = "integer"
			tf = "uint64"
		case "uint32":
			fallthrough
		case "int64":
			t = "integer"
			tf = "int64"
		case "int32":
			t = "integer"
			tf = "int32"
		case "bool":
			t = "Boolean"
		default:
			t = fmt.Sprintf("%s.%s", c.PackageName, field.Type)
		}

		//	Is it repeated?
		if field.Repeated {
			t = fmt.Sprintf("[%s!]", t)
		}

		//	Is it nullable?
		if !isNullable(field.Field) {
			t = fmt.Sprintf("%s!", t)
		}

		//	Get comment information
		//out += genComment(field.Comment)

	}

	return retval
}
