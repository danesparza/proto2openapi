package cmd

import (
	"github.com/danesparza/proto2openapi/internal/converter"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/spf13/cobra"
)

var ProtoPath string
var OutputPath string

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
	generateCmd.Flags().StringVarP(&OutputPath, "out", "o", "./apischema.yaml", "OpenAPI schema file to write to")
}

func generate(cmd *cobra.Command, args []string) {

	//	Open the proto path
	reader, err := os.Open(ProtoPath)
	defer reader.Close()
	if err != nil {
		log.Err(err).Str("proto", ProtoPath).Msg("Problem opening the proto file")
		return
	}

	//	Create a new Converter type:
	c := converter.Converter{
		Source: reader,
	}

	//	Generate YAML docs
	retval, err := c.ConvertToYAML()
	if err != nil {
		log.Err(err).Str("proto", ProtoPath).Msg("Problem generating docs")
	}

	//	Spit out the docs:
	err = os.WriteFile(OutputPath, []byte(retval), 0666)
	if err != nil {
		log.Err(err).Str("proto", ProtoPath).Str("outfile", OutputPath).Msg("Problem writing to output file")
	} else {
		log.Info().Str("proto", ProtoPath).Str("outfile", OutputPath).Msg("Finished generating docs to output file")
	}

}
