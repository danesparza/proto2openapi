package cmd

import (
	"github.com/emicklei/proto"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/spf13/cobra"
)

var ProtoPath string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate OpenAPI documentation file from the protobuf file",
	Long:  `Generate OpenAPI documentation file from the protobuf file`,
	Run:   generate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&ProtoPath, "proto", "p", "", "Source protobuf file to read from")
}

func generate(cmd *cobra.Command, args []string) {
	log.Trace().Msg("Generate called")

	reader, err := os.Open(ProtoPath)
	defer reader.Close()

	if err != nil {
		log.Err(err).Str("proto", ProtoPath).Msg("Problem opening the proto file")
		return
	}

	// Parse
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Err(err).Str("proto", ProtoPath).Msg("Problem parsing the proto file")
		return
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

	packageName := pkg.Name
	log.Info().Str("proto", ProtoPath).Str("package", packageName).Msg("Parsed proto file")

}
