/*
Copyright © 2023 Dataflows
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/thedataflows/go-commons/pkg/config"
	"github.com/thedataflows/namecheap-cli/pkg/constants"

	"github.com/spf13/cobra"
)

type requestParameters struct {
	sandbox  string
	apiKey   string
	username string
	tld      string
	sld      string
	clientIP string
}

const (
	keyCommonSandbox  = "sandbox"
	keyCommonApiKey   = "key"
	keyCommonUsername = "username"
	keyCommonTld      = "tld"
	keyCommonSld      = "sld"
	keyCommonClientIp = "client-ip"
	namecheapApiUrl   = "https://api.%snamecheap.com/xml.response?apiuser=%s&apikey=%s&username=%s&SLD=%s&TLD=%s&ClientIP=%s&Command=namecheap.domains.dns.%s"
)

var (
	supportedFormats = []string{"xml", "yaml", "json"}

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "namecheap-cli",
		Short: "Namecheap DNS command line interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Long = fmt.Sprintf(
				"%s\n\nAll flags values can be provided via env vars starting with %s_*\nTo pass a subcommand (e.g. 'serve') flag, use %s_GET_FLAGNAME=somevalue",
				cmd.Short,
				configOpts.EnvPrefix,
				configOpts.EnvPrefix,
			)
			_ = cmd.Help()
		},
	}

	configOpts = config.DefaultConfigOpts(
		&config.Opts{
			EnvPrefix: constants.ViperEnvPrefix,
		},
	)
)

func setCommonParameters(cmd *cobra.Command) *requestParameters {
	s := ""
	if config.ViperGetBool(cmd, keyCommonSandbox) {
		s = "sandbox."
	}
	return &requestParameters{
		sandbox:  s,
		apiKey:   config.ViperGetString(cmd, keyCommonApiKey),
		username: config.ViperGetString(cmd, keyCommonUsername),
		sld:      config.ViperGetString(cmd, keyCommonSld),
		tld:      config.ViperGetString(cmd, keyCommonTld),
		clientIP: config.ViperGetString(cmd, keyCommonClientIp),
	}
}

func initConfig() {
	config.InitConfig(configOpts)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().AddFlagSet(configOpts.Flags)
	config.ViperBindPFlagSet(rootCmd, configOpts.Flags)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
