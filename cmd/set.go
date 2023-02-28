/*
Copyright © 2023 Dataflows
*/
package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/thedataflows/go-commons/pkg/config"
	"github.com/thedataflows/go-commons/pkg/log"
	"github.com/thedataflows/go-commons/pkg/stringutil"
	"github.com/thedataflows/namecheap-cli/pkg/namecheap"
	"k8s.io/utils/strings/slices"

	"github.com/spf13/cobra"
)

const (
	keySetInputFile   = "input-file"
	keySetInputFormat = "input-format"
	keySetTimeout     = "timeout"
)

var (
	requiredSetFlags = []string{keyCommonApiKey, keyCommonUsername}

	setCmd = &cobra.Command{
		Use:     "set",
		Short:   "Upload Namecheap DNS configuration",
		Long:    ``,
		Aliases: []string{"s"},
		Run:     RunSet,
	}
)

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().Bool(keyCommonSandbox, false, "Use Namecheap sandbox API")
	setCmd.Flags().StringP(keyCommonApiKey, "k", "", "[Required] Namecheap API key")
	setCmd.Flags().StringP(keyCommonUsername, "u", "", "[Required] Namecheap user")
	setCmd.Flags().StringP(keyCommonTld, "t", "", "Namecheap top-level domain, e.g.: 'com'. Can be read from the input file")
	setCmd.Flags().StringP(keyCommonSld, "s", "", "Namecheap second-level domain, e.g.: 'example'. Can be read from the input file")
	setCmd.Flags().String(keyCommonClientIp, "127.0.0.1", "Client IP. This is not really required")

	setCmd.Flags().StringP(keySetInputFile, "i", "", "Input file. If omitted, stdin is used until 2 consecutive newlines are detected")
	setCmd.Flags().String(keySetInputFormat, supportedFormats[0], fmt.Sprintf("Input format. Supported: %v", supportedFormats))
	setCmd.Flags().Duration(keySetTimeout, 10, "Request timeout")

	config.ViperBindPFlagSet(setCmd, nil)
}

// RunSet uploads the Namecheap DNS configuration
func RunSet(cmd *cobra.Command, args []string) {
	// Validations
	config.CheckRequiredFlags(cmd, requiredSetFlags)

	format := config.ViperGetString(cmd, keySetInputFormat)
	if !slices.Contains(supportedFormats, format) {
		log.Fatalf("Input format '%s' is not supported. Please use one of: %v", format, supportedFormats)
	}

	input := unmarshal(
		format,
		readInput(cmd),
	)
	// try to get tld and sld from input data
	domainSegments := strings.Split(input.CommandResponse.DomainDNSGetHostsResult.Domain, ".")
	if len(config.ViperGetString(cmd, keyCommonSld)) == 0 {
		if len(domainSegments) < 2 || len(domainSegments[0]) == 0 {
			log.Fatalf("Neither --%s was specified nor '/ApiResponse/CommandResponse/DomainDNSGetHostsResult/@Domain' was set in the input!", keyCommonSld)
		}
		config.ViperSet(cmd, keyCommonSld, domainSegments[0])
	}
	if len(config.ViperGetString(cmd, keyCommonTld)) == 0 {
		if len(domainSegments) < 2 || len(domainSegments[1]) == 0 {
			log.Fatalf("Neither --%s was specified nor CommandResponse.DomainDNSGetHostsResult.Domain was set in the input!", keyCommonTld)
		}
		config.ViperSet(cmd, keyCommonTld, domainSegments[1])
	}

	upload(cmd, input, config.ViperGetDuration(cmd, keySetTimeout))
}

// upload performs a POST request on the Namecheap API endpoint with the
func upload(cmd *cobra.Command, input *namecheap.ApiResponse, timeout time.Duration) {
	parentReqParams := setCommonParameters(cmd)

	log.Info("Uploading Namecheap DNS configuration")

	url := fmt.Sprintf(
		namecheapApiUrl,
		parentReqParams.sandbox,
		parentReqParams.username,
		parentReqParams.apiKey,
		parentReqParams.username,
		parentReqParams.sld,
		parentReqParams.tld,
		parentReqParams.clientIp,
		"setHosts",
	)
	log.Debug(url)

	requestBody := ""
	for _, host := range input.CommandResponse.DomainDNSGetHostsResult.Host {
		if len(host.HostId) > 0 {
			requestBody = stringutil.ConcatStrings(
				requestBody,
				"HostName", host.HostId, "=", host.Name, "&",
				"RecordType", host.HostId, "=", host.Type, "&",
				"Address", host.HostId, "=", host.Address, "&",
				"MXPref", host.HostId, "=", host.MXPref, "&",
				"TTL", host.HostId, "=", host.TTL, "&",
				"FriendlyName", host.HostId, "=", host.FriendlyName, "&",
				"IsActive", host.HostId, "=", host.IsActive, "&",
			)
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: time.Second * timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}

	response := &namecheap.ApiResponse{}
	err = xml.Unmarshal(body, response)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if response.Status != "OK" {
		messages := ""
		for _, e := range response.Errors.Error {
			messages = fmt.Sprintf("%s%s: %s", messages, e.Number, e.Text)
		}
		log.Fatalf("Received errors from the api server: \n%v", messages)
	}

	log.Debugf("Raw response: \n%s", string(body))
	log.Infof("Success. Execution time: %s", response.ExecutionTime)
}
