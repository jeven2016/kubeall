package main

import (
	"embed"
	_ "embed"
	"github.com/spf13/cobra"
	"kubeall.io/api-server/pkg/server"
	"kubeall.io/api-server/pkg/types"
	"log"
)

// internal configuration file
//
//go:embed resources/internal_conf.yaml
var internalCfg string

//go:embed resources/banner.txt
var banner string

// i18n resources
//
//go:embed resources/locales/*
var localeFs embed.FS

// the external configuration file to be used for overriding the existing configuration
const flagName = "config"

// system initializing
// run command： auth-proxy -c external_config.yaml
func main() {
	//显示banner信息
	println(banner)

	var rootCmd = &cobra.Command{
		Version: "0.1.0",
		Use:     "api-server",
		Short:   "api server",
		Run:     startServer,
	}

	rootCmd.Flags().StringP(flagName, "c", "", "the absolute path of yaml config file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// startServer initializes and starts the server with the provided configuration.
// It retrieves the custom configuration file path from command-line flags,
// sets up the startup parameters, and registers the server modules.
//
// Parameters:
//   - cmd: A pointer to cobra.Command representing the command being run.
//   - args: A slice of strings containing any additional command-line arguments.
//
// The function does not return any value. If an error occurs while retrieving
// the custom configuration file path, it prints an error message and returns early.
func startServer(cmd *cobra.Command, args []string) {
	customCfg, err := cmd.Flags().GetString(flagName)
	if err != nil {
		println("error occurs: %s", err.Error())
		return
	}

	params := &types.StartupParams{
		InternalConfig:   []byte(internalCfg),
		CustomConfigPath: customCfg,
		DefaultCfgPaths:  []string{},
	}
	server.RegisterModules(params, localeFs)
}
