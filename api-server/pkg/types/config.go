package types

type StartupParams struct {
	InternalConfig   []byte
	CustomConfigPath string
	DefaultCfgPaths  []string
}

type Config interface {
	Validate() error
	Complete() error
}

type LogConfig struct {
	Enabled         bool   `koanf:"enabled"`
	LogLevel        string `koanf:"logLevel"`
	LogPath         string `koanf:"logPath"`
	OutputConsole   bool   `koanf:"outputToConsole"`
	FileName        string `koanf:"fileName"`
	MaxSizeInMB     int    `koanf:"maxSizeInMB"`
	MaxAgeInDays    int    `koanf:"maxAgeInDays"`
	MaxBackups      int    `koanf:"maxBackups"`
	Compress        bool   `koanf:"compress"`
	PrintErrorStack bool   `koanf:"printErrorStack"`
}

// HttpSetting 服务
type HttpSetting struct {
	Address string `koanf:"address" yaml:"address"`
	Port    uint   `koanf:"port" yaml:"port"`
}

type StorageClassConfig struct {
	NumberOfReplicas     string `koanf:"numberOfReplicas"`
	StaleReplicaTimeout  string `koanf:"staleReplicaTimeout"`
	AllowVolumeExpansion bool   `koanf:"allowVolumeExpansion"`
	ReclaimPolicy        string `koanf:"reclaimPolicy"`
	Provisioner          string `koanf:"provisioner"`
}

type ServerConfig struct {
	ApplicationName    string              `koanf:"applicationName"`
	LogSetting         *LogConfig          `koanf:"logConfig"`
	Http               *HttpSetting        `koanf:"http" yaml:"http"`
	KubeConfig         string              `koanf:"kubeConfig"`
	StorageClassConfig *StorageClassConfig `koanf:"storeClass"`
}

func (s ServerConfig) GetServerConfig() *ServerConfig {
	return &s
}

func (s ServerConfig) Validate() error {
	return nil
}

func (s ServerConfig) Complete() error {
	return nil
}
