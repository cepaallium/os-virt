package options

import (
	"fmt"
	"github.com/spf13/viper"
	"os-virt/pkg/utils/log"
)

const (
	defaultConfiguration string = "config/os-virt.yaml"
)

type VirtConfig struct {
	WebhookPort   int `mapstructure:"web-hook-port"`
	WebserverPort int `mapstructure:"web-server-port"`

	LogFile    string `mapstructure:"log-file"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxBackups int    `mapstructure:"max-backups"`
	MaxAge     int    `mapstructure:"max-age"`
	Compress   bool   `mapstructure:"compress"`

	LogLevel        string `mapstructure:"log-level"`
	JsonEncode      bool   `mapstructure:"json-encode"`
	StacktraceLevel string `mapstructure:"stacktrace-level"`
}

func (mc *VirtConfig) LoadConfig() error {
	v := viper.New()
	v.SetConfigFile(defaultConfiguration)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return err
		} else {
			return fmt.Errorf("error parsing configuration file %s", err)
		}
	}
	if err := v.Unmarshal(mc); err != nil {
		return err
	}
	mc.LoadEnvValues()
	//mc.LoadFlagValues()
	return nil
}

func (c *VirtConfig) LoadEnvValues() {

}

func (vc *VirtConfig) MakeLoggerCfg() *log.Config {
	loggerConfig := log.Config{}
	loggerConfig.LogFile = vc.LogFile
	loggerConfig.LogLevel = vc.LogLevel
	loggerConfig.Compress = vc.Compress
	loggerConfig.MaxAge = vc.MaxAge
	loggerConfig.MaxSize = vc.MaxSize
	loggerConfig.MaxBackups = vc.MaxBackups
	loggerConfig.JsonEncode = vc.JsonEncode
	loggerConfig.StacktraceLevel = vc.StacktraceLevel
	return &loggerConfig
}
