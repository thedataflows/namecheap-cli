/*
Copyright © 2023 Dataflows
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/thedataflows/go-commons/pkg/config"
	"github.com/thedataflows/go-commons/pkg/file"
	"github.com/thedataflows/go-commons/pkg/log"
	"github.com/thedataflows/namecheap-cli/pkg/namecheap"
	"gopkg.in/yaml.v3"
	"k8s.io/utils/strings/slices"

	"github.com/spf13/cobra"
)

const keyConvertForce = "force"

var (
	convertCmd = &cobra.Command{
		Use:     "convert",
		Short:   "Convert Namecheap DNS configuration between local storage formats",
		Long:    ``,
		Aliases: []string{"c"},
		Run:     RunConvert,
	}
)

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringP(keyCommonTld, "t", "", "Namecheap top-level domain, e.g.: 'com'")
	convertCmd.Flags().StringP(keyCommonSld, "s", "", "Namecheap second-level domain, e.g.: 'example'")

	convertCmd.Flags().StringP(keySetInputFile, "i", "", "Input file. If omitted, stdin is used until 2 consecutive newlines are detected")
	convertCmd.Flags().StringP(keyGetOutputFile, "o", "", "Output file. If omitted, outputs to stdout")
	convertCmd.Flags().String(keySetInputFormat, supportedFormats[0], fmt.Sprintf("Input format. Supported: %v", supportedFormats))
	convertCmd.Flags().String(keyGetOutputFormat, supportedFormats[1], fmt.Sprintf("Output format. Supported: %v", supportedFormats))
	convertCmd.Flags().Bool(keyConvertForce, false, "Overwrite the file if exists")

	config.ViperBindPFlagSet(convertCmd, nil)
}

// RunConvert converts between one supported format to another
func RunConvert(cmd *cobra.Command, args []string) {
	// Validations
	inputFormat := config.ViperGetString(cmd, keySetInputFormat)
	if !slices.Contains(supportedFormats, inputFormat) {
		log.Fatalf("Input format '%s' is not supported. Please use one of: %v", inputFormat, supportedFormats)
	}
	outputFormat := config.ViperGetString(cmd, keyGetOutputFormat)
	if !slices.Contains(supportedFormats, outputFormat) {
		log.Fatalf("Output format '%s' is not supported. Please use one of: %v", outputFormat, supportedFormats)
	}
	if strings.EqualFold(inputFormat, outputFormat) {
		log.Fatalf("Input format is the same as output format, they must be different")
	}

	inputData := readInput(cmd)
	input := unmarshal(inputFormat, inputData)
	sld := config.ViperGetString(cmd, keyCommonSld)
	tld := config.ViperGetString(cmd, keyCommonTld)
	if len(sld) > 0 && len(tld) > 0 {
		input.CommandResponse.DomainDNSGetHostsResult.Domain = fmt.Sprintf("%s.%s", sld, tld)
	}

	output := marshal(outputFormat, input)
	writeOutput(cmd, output)
}

// readInput reads data from stdin or file, if provided
func readInput(cmd *cobra.Command) *[]byte {
	inputHandle := os.Stdin
	inputFileName := config.ViperGetString(cmd, keySetInputFile)
	if len(inputFileName) > 0 {
		if !file.IsFile(inputFileName) {
			log.Fatalf("'%s' is not accessible", inputFileName)
		}
		var err error
		inputHandle, err = os.Open(inputFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer inputHandle.Close()
	}

	var inputData []byte
	scanner := bufio.NewScanner(inputHandle)
	enterPressed := 0
	for scanner.Scan() {
		if inputHandle == os.Stdin {
			if scanner.Text() == "" {
				enterPressed++
			}
			if enterPressed > 1 {
				break
			}
		}
		inputData = append(inputData, scanner.Bytes()...)
		inputData = append(inputData, '\n')
	}
	return &inputData
}

// unmarshal populates namecheap.ApiResponse
func unmarshal(inputFormat string, input *[]byte) *namecheap.ApiResponse {
	var (
		inputMarshalled = &namecheap.ApiResponse{}
		err             error
	)
	switch inputFormat {
	case supportedFormats[0]:
		err = xml.Unmarshal(*input, inputMarshalled)
	case supportedFormats[1]:
		err = yaml.Unmarshal(*input, inputMarshalled)
	case supportedFormats[2]:
		err = json.Unmarshal(*input, inputMarshalled)
	}
	if err != nil {
		log.Fatalf("Failed to unmarshal: %s", err)
	}
	return inputMarshalled
}

// marshal is marshaling namecheap.ApiResponse to the specified format
func marshal(format string, apiresponse *namecheap.ApiResponse) *[]byte {
	var (
		output []byte
		err    error
	)
	switch format {
	case supportedFormats[0]:
		output, err = xml.MarshalIndent(apiresponse, "", "  ")
	case supportedFormats[1]:
		output, err = yaml.Marshal(apiresponse)
	case supportedFormats[2]:
		output, err = json.MarshalIndent(apiresponse, "", "  ")
	}
	if err != nil {
		log.Fatalf("Failed to marshal format '%s': %s", format, err)
	}
	return &output
}

// writeOutput writes to a specified file or stdout
func writeOutput(cmd *cobra.Command, output *[]byte) {
	outputFileName := config.ViperGetString(cmd, keyGetOutputFile)
	if file.IsFile(outputFileName) && !config.ViperGetBool(cmd, keyConvertForce) {
		log.Fatalf("'%s' exists, but without the --%s flag, will not overwrite it!", outputFileName, keyConvertForce)
	}

	if len(outputFileName) > 0 {
		destination, err := os.Create(outputFileName)
		if err != nil {
			log.Fatalf("Failed to create file '%s' because: %s", outputFileName, err)
		}
		defer destination.Close()

		if _, err := destination.Write(*output); err != nil {
			log.Fatalf("Failed to write to file '%s' because: %s", outputFileName, err)
		}
	} else {
		fmt.Printf("%s\n", *output)
	}
}
