package converter

import (
	"fmt"
	"github.com/emicklei/proto"
	"github.com/rs/zerolog/log"
	"strings"
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
	return fmt.Sprint("  schemas:\n")
}

// GetYAMLTypeFromMessage emits OpenAPI YAML type information for a proto message
func (c *Converter) GetYAMLTypeFromMessage(message *proto.Message) string {
	retval := ""

	//	Write the basic type information
	retval += fmt.Sprintf("    %s.%s:\n      type: object\n", c.PackageName, message.Name)

	//	If this has a message-level comment, use it
	if message.Comment != nil {
		retval += fmt.Sprintf("      description: %s", c.GetYAMLComment(message.Comment, 8))
	}

	//	If this has fields, write 'properties' and then write the fields under that.
	if len(message.Elements) > 0 {
		retval += fmt.Sprintf("      properties:\n")
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
		retval += fmt.Sprintf("        %s:\n", name)

		//	Map the proto type to the OpenAPI type
		//	Protobuf type reference: https://protobuf.dev/programming-guides/proto3/#scalar
		//	OpenAPI type reference: https://swagger.io/docs/specification/data-models/data-types/
		var t string
		var tf string
		customType := false
		switch field.Type {
		case "double":
			t = "number"
			tf = "double"
		case "float":
			t = "number"
			tf = "float"
		case "int32":
			t = "integer"
			tf = "int32"
		case "int64":
			t = "integer"
			tf = "int64"
		case "uint32":
			t = "number"
		case "uint64":
			t = "number"
		case "sint32":
			t = "integer"
			tf = "int32"
		case "sint64":
			t = "integer"
			tf = "int64"
		case "fixed32":
			t = "number"
		case "fixed64":
			t = "number"
		case "sfixed32":
			t = "integer"
			tf = "int32"
		case "sfixed64":
			t = "integer"
			tf = "int64"
		case "bool":
			t = "boolean"
		case "string":
			t = "string"
		case "bytes":
			/* Not sure the best way to represent this in OpenAPI types */
			t = "string"
		default:
			customType = true
			t = fmt.Sprintf("%s.%s", c.PackageName, field.Type)
		}

		//	FORMAT THE FIELD INFORMATION

		//	It's not an array
		if !field.Repeated {

			//	If it's a custom type, use ref formatting
			if customType {
				retval += fmt.Sprintf("          $ref: '#/components/schemas/%s'\n", t)
			}

			//	If it's not a custom type, use regular formatting
			if !customType {
				retval += fmt.Sprintf("          type: %s\n", t)

				//	If we have a type format, show it
				if tf != "" {
					retval += fmt.Sprintf("          format: %s\n", tf)
				}
			}

		}

		//	It is an array
		if field.Repeated {
			//	If it's a custom type, use ref formatting
			if customType {
				retval += fmt.Sprintf("          type: array\n          items:\n            $ref: '#/components/schemas/%s'\n", t)
			}

			//	If it's not a custom type, use regular formatting
			if !customType {
				retval += fmt.Sprintf("          type: array\n          items:\n            type: %s\n", t)

				//	If we have a type format, show it
				if tf != "" {
					retval += fmt.Sprintf("            format: %s\n", tf)
				}
			}
		}

		//	Get field-level comment information
		if field.Comment != nil {
			retval += fmt.Sprintf("          description: %s", c.GetYAMLComment(field.Comment, 12))
		}

		if field.InlineComment != nil {
			retval += fmt.Sprintf("          description: %s", c.GetYAMLComment(field.InlineComment, 12))
		}

	}

	return retval
}

// GetYAMLComment returns a properly formatted YAML comment
// for the given proto comment (or an empty string if the comment is nil)
func (c *Converter) GetYAMLComment(comment *proto.Comment, indentLvl int) string {
	retval := ""

	if comment == nil {
		return retval
	}

	indent := strings.Repeat(" ", indentLvl)

	//	A multi-line comment
	if len(comment.Lines) > 1 {
		retval += "|-\n"

		for _, line := range comment.Lines {
			if len(strings.TrimSpace(line)) > 0 {
				retval += indent
				retval += strings.TrimSpace(line)
				retval += "\n"
			}
		}

	} else {
		//	A single line comment
		retval = fmt.Sprintf("%s\n", strings.TrimSpace(comment.Message()))
	}

	return retval
}
