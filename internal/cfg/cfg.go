package cfg

import (
	"fmt"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

const SettingsName string = "settings"
const SettingsFile string = "settings.yaml"

type Config struct {
	Sync struct {
		Server   string // `mapstructure:server`
		User     string // `mapstructure:user`
		Password string // `mapstructure:password`
		Port     string // `mapstructure:port`
		Folder   string // `mapstructure:folder`
		Accounts string // `mapstructure:accounts`
	} // `mapstructure:sync`
	Global struct {
		Logging      string // `mapstructure:logging`
		Account      string // `mapstructure:account`
		Local_folder string // `mapstructure:local_folder`
	} // `mapstructure:global`
	Import struct {
		Remote_limit int // `mapstrucutre:"remote_limit"`
		Local_limit  int // `mapstrucutre:"local_limit"`
	} // `mapstructure:"import"`
}

func (c *Config) Write() error {
	settings := viper.AllSettings()
	bs, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(SettingsFile, bs, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) New() (*Config, error) {
	fmt.Println("Loading config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigName(SettingsName)
	viper.SetDefault("Sync.Folder", "sync")
	viper.SetDefault("Sync.Accounts", "all")
	viper.SetDefault("Sync.Server", "sftp-server")
	viper.SetDefault("Sync.User", "sftp-user")
	viper.SetDefault("Sync.Password", "sftp-password")
	viper.SetDefault("Sync.Port", 22)
	viper.SetDefault("Global.Logging", "INFO")
	viper.SetDefault("Global.Local_folder", "APPDATA")
	viper.SetDefault("Import.Remote_limit", 1)
	viper.SetDefault("Import.Local_limit", 10)
	err := viper.ReadInConfig()
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			c.Write()
		}
		return c, err
	}
	err = viper.Unmarshal(c)
	if err != nil {
		return c, err
	}
	return c, nil
}
