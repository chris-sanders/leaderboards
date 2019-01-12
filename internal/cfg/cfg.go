package cfg

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const SettingsName string = "settings"
const SettingsFile string = "settings.yaml"

type Config struct {
	Sync struct {
		Server   string
		User     string
		Password string
		Port     string
		Folder   string
		Accounts string
	}
	Global struct {
		Logging      string
		Account      string
		Local_folder string
		Local_file   string `json:"-"`
	}
	Import struct {
		Remote_limit int
		Local_limit  int
	}
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

func postprocess(c *Config) {
	var full_path string
	switch c.Global.Local_folder {
	case "APPDATA":
		full_path = filepath.Join(os.ExpandEnv("${userprofile}/AppData/LocalLow/Hyperbolic Magnetism/Beat Saber"), "LocalLeaderboards.dat")
	default:
		full_path = filepath.Join(c.Global.Local_folder, "LocalLeaderboards.dat")
	}
	c.Global.Local_file = full_path
}

func (c *Config) New() (*Config, error) {
	log.Debug("Loading config")
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
	viper.SetDefault("Global.Account", "unique-name")
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
	postprocess(c)
	return c, nil
}
