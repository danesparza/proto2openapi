package converter

import (
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
