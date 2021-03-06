package conf

import (
	"github.com/yangxikun/multiconfig"
)

type ConfigTOML struct {
	Database DatabaseConfig
	Server   ServerConfig
	Auth     struct {
		Enable   bool
		Secret   string
		Accounts map[string]string
		Skips    []string
	}

	Log struct {
		FilePath string
		FileName string
	}
}

var Config *ConfigTOML
var ConfigTomlPath string
var ConfigArgs string

func Init(tomlPath string, args ...string) *ConfigTOML {
	Config = &ConfigTOML{}
	ConfigTomlPath = tomlPath
	if len(args) > 0 {
		ConfigArgs = args[0]
	}
	loader := multiconfig.NewWithPath(tomlPath)
	loader.MustLoad(Config)
	return Config
}
