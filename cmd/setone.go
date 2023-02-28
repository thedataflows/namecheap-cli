/*
Copyright © 2023 Dataflows
*/
package cmd

import (
	"github.com/thedataflows/go-commons/pkg/config"
	"github.com/thedataflows/go-commons/pkg/log"
	"github.com/thedataflows/namecheap-cli/pkg/namecheap"

	"github.com/spf13/cobra"
)

const (
	setOneKeyName         = "name"
	setOneKeyType         = "type"
	setOneKeyAddress      = "address"
	setOneKeyMXPref       = "mxpref"
	setOneKeyTTL          = "ttl"
	setOneKeyFriendlyName = "friendlyname"
	setOneKeyIsActive     = "isactive"
	setOneKeyDelete       = "delete"
)

var (
	requiredSetOneFlags = []string{keyCommonApiKey, keyCommonUsername, keyCommonTld, keyCommonSld, setOneKeyName, setOneKeyType, setOneKeyAddress}

	setOneCmd = &cobra.Command{
		Use:     "setone",
		Short:   "Create/update/delete a single DNS entry",
		Long:    ``,
		Aliases: []string{"o"},
		Run:     RunSetOne,
	}
)

func init() {
	rootCmd.AddCommand(setOneCmd)

	setOneCmd.Flags().Bool(keyCommonSandbox, false, "Use Namecheap sandbox API")
	setOneCmd.Flags().StringP(keyCommonApiKey, "k", "", "[Required] Namecheap API key")
	setOneCmd.Flags().StringP(keyCommonUsername, "u", "", "[Required] Namecheap user")
	setOneCmd.Flags().StringP(keyCommonTld, "t", "", "[Required] Namecheap top-level domain, e.g.: 'com'")
	setOneCmd.Flags().StringP(keyCommonSld, "s", "", "[Required] Namecheap second-level domain, e.g.: 'example'")
	setOneCmd.Flags().String(keyCommonClientIp, "127.0.0.1", "Client IP. This is not really required")

	setOneCmd.Flags().String(setOneKeyName, "", "[Required] Record name")
	setOneCmd.Flags().String(setOneKeyType, "", "[Required] Record type")
	setOneCmd.Flags().String(setOneKeyAddress, "", "[Required] Record value")
	setOneCmd.Flags().String(setOneKeyMXPref, "", "MXPref")
	setOneCmd.Flags().String(setOneKeyTTL, "1799", "Time to live in seconds. 1799 is Namecheap's equivalent to 'Automatic'")
	setOneCmd.Flags().String(setOneKeyFriendlyName, "", "Friendly name")
	setOneCmd.Flags().Bool(setOneKeyIsActive, true, "Active state")

	setOneCmd.Flags().Bool(setOneKeyDelete, false, "Delete DNS entry")

	setOneCmd.Flags().Duration(keyGetTimeout, 10, "Request timeout")

	config.ViperBindPFlagSet(setOneCmd, nil)
}

// RunSetOne downloads the Namecheap DNS configuration first, merges single entry from input then uploads the entire configuration back
func RunSetOne(cmd *cobra.Command, args []string) {
	// Validations
	config.CheckRequiredFlags(cmd, requiredSetOneFlags)

	inputHost := &namecheap.Host{
		HostId:       "1",
		Name:         config.ViperGetString(cmd, setOneKeyName),
		Type:         config.ViperGetString(cmd, setOneKeyType),
		Address:      config.ViperGetString(cmd, setOneKeyAddress),
		MXPref:       config.ViperGetString(cmd, setOneKeyMXPref),
		TTL:          config.ViperGetString(cmd, setOneKeyTTL),
		FriendlyName: config.ViperGetString(cmd, setOneKeyFriendlyName),
		IsActive:     config.ViperGetString(cmd, setOneKeyIsActive),
	}

	timeout := config.ViperGetDuration(cmd, keyGetTimeout)
	delete := config.ViperGetBool(cmd, setOneKeyDelete)

	// download current DNS configuration
	apiresponse := download(cmd, timeout)

	found := false
	for i, host := range apiresponse.CommandResponse.DomainDNSGetHostsResult.Host {
		if host.Name == inputHost.Name && host.Type == inputHost.Type {
			log.Debugf("Matched host: %#v", host)

			if delete {
				apiresponse.CommandResponse.DomainDNSGetHostsResult.Host[i] = namecheap.Host{}
				break
			}

			host.Address = inputHost.Address
			if len(inputHost.MXPref) > 0 {
				host.MXPref = inputHost.MXPref
			}
			if len(inputHost.TTL) > 0 {
				host.TTL = inputHost.TTL
			}
			if len(inputHost.FriendlyName) > 0 {
				host.FriendlyName = inputHost.FriendlyName
			}
			host.IsActive = inputHost.IsActive
			apiresponse.CommandResponse.DomainDNSGetHostsResult.Host[i] = host

			found = true
			break
		}
	}

	if !found && !delete {
		apiresponse.CommandResponse.DomainDNSGetHostsResult.Host = append(
			apiresponse.CommandResponse.DomainDNSGetHostsResult.Host,
			*inputHost,
		)
	}

	// upload new DNS configuration
	upload(cmd, apiresponse, timeout)
}
