package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"kubeall.io/api-server/pkg/infra/utils"
	"kubeall.io/api-server/pkg/types"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func NewServerConfig(params *types.StartupParams) (types.Config, error) {
	cfg := &types.ServerConfig{}

	if err := loadConfig(params.InternalConfig, cfg, &params.CustomConfigPath, params.DefaultCfgPaths); err != nil {
		println("error occurs while loading config file: %s", err.Error())
		return nil, err
	}
	return cfg, nil
}

// loadConfig loads the configuration files
func loadConfig(internalCfg []byte, config types.Config, extraConfigFilePath *string, defaultCfgPaths []string) error {

	//load internal config
	if internalCfg != nil {
		if err := k.Load(rawbytes.Provider(internalCfg), yaml.Parser()); err != nil {
			return err
		}
	}

	cfgPaths := defaultCfgPaths
	if cfgPaths == nil {
		cfgPaths = []string{}
	}
	if extraConfigFilePath != nil {
		cfgPaths = append(cfgPaths, *extraConfigFilePath)
	}

	// load external configs
	for _, f := range cfgPaths {
		if exists, err := utils.IsFileExists(f); err != nil {
			continue
		} else if exists {
			if err = k.Load(file.Provider(f), yaml.Parser()); err != nil {
				return err
			}
		}
	}

	if err := k.Unmarshal("", config); err != nil {
		return err
	}

	return nil
}
