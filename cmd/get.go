/*
Copyright © 2023 Dataflows
*/
package cmd

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thedataflows/go-commons/pkg/config"
	"github.com/thedataflows/go-commons/pkg/log"
	"github.com/thedataflows/namecheap-cli/pkg/namecheap"
	"k8s.io/utils/strings/slices"

	"github.com/spf13/cobra"
)

const (
	keyGetOutputFile   = "output-file"
	keyGetOutputFormat = "output-format"
	keyGetTimeout      = "timeout"
)

var (
	requiredGetFlags = []string{keyCommonApiKey, keyCommonUsername, keyCommonTld, keyCommonSld}

	getCmd = &cobra.Command{
		Use:     "get",
		Short:   "Download Namecheap DNS configuration",
		Long:    ``,
		Aliases: []string{"g"},
		Run:     RunGet,
	}
)

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().Bool(keyCommonSandbox, false, "Use Namecheap sandbox API")
	getCmd.Flags().StringP(keyCommonApiKey, "k", "", "[Required] Namecheap API key")
	getCmd.Flags().StringP(keyCommonUsername, "u", "", "[Required] Namecheap user")
	getCmd.Flags().StringP(keyCommonTld, "t", "", "[Required] Namecheap top-level domain, e.g.: 'com'")
	getCmd.Flags().StringP(keyCommonSld, "s", "", "[Required] Namecheap second-level domain, e.g.: 'example'")
	getCmd.Flags().String(keyCommonClientIp, "127.0.0.1", "Client IP. This is not really required")

	getCmd.Flags().StringP(keyGetOutputFile, "o", "", "Output file. If omitted, outputs to stdout")
	getCmd.Flags().String(keyGetOutputFormat, supportedFormats[0], fmt.Sprintf("Output format. Supported: %v", supportedFormats))
	getCmd.Flags().Bool(keyConvertForce, false, "Force overwriting the file if exists")
	getCmd.Flags().Duration(keyGetTimeout, 10, "Request timeout")

	config.ViperBindPFlagSet(getCmd, nil)
}

// RunGet downloads the Namecheap DNS configuration and saves it as specified format
func RunGet(cmd *cobra.Command, args []string) {
	config.CheckRequiredFlags(cmd, requiredGetFlags)

	format := config.ViperGetString(cmd, keyGetOutputFormat)
	if !slices.Contains(supportedFormats, format) {
		log.Fatalf("Output format '%s' is not supported. Please use one of: %v", format, supportedFormats)
	}

	output := marshal(
		format,
		download(
			cmd,
			config.ViperGetDuration(cmd, keyGetTimeout),
		),
	)
	writeOutput(cmd, output)
}

// download performs a GET request on the Namecheap API endpoint returning the response body unmarshaled from XML
func download(cmd *cobra.Command, timeout time.Duration) *namecheap.ApiResponse {
	parentReqParams := setCommonParameters(cmd)

	log.Info("Downloading Namecheap DNS configuration")

	url := fmt.Sprintf(
		namecheapApiUrl,
		parentReqParams.sandbox,
		parentReqParams.username,
		parentReqParams.apiKey,
		parentReqParams.username,
		parentReqParams.sld,
		parentReqParams.tld,
		parentReqParams.clientIP,
		"getHosts",
	)
	log.Debug(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	req.Header.Set("Cache-Control", "no-cache")

	// init & call
	client := &http.Client{
		Timeout: time.Second * timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	defer resp.Body.Close()

	// read
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}

	// unmarshal
	response := &namecheap.ApiResponse{}
	err = xml.Unmarshal(body, response)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// status check
	if response.Status != "OK" {
		messages := ""
		for _, e := range response.Errors.Error {
			messages = fmt.Sprintf("%s%s: %s\n", messages, e.Number, e.Text)
		}
		log.Fatalf("Received errors from the api server: \n%v", messages)
	}

	log.Debugf("Raw response: \n%s", string(body))
	log.Infof("Success. Execution time: %s", response.ExecutionTime)

	return response
}
