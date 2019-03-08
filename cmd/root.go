package cmd

import (
	"strings"

	"github.com/0987363/viper"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "config",
	Short: "config backend server",
}

var configFilePath string

func init() {
	// Register a global option for config file path
	RootCmd.PersistentFlags().StringVarP(
		&configFilePath, "config", "c", "", "Path to the config file",
	)
}

// LoadConfiguration should be used as PreRun function by commands that need
// get content from configuration file.
func LoadConfiguration(cmd *cobra.Command, args []string) {
	// Load the configuration
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/druid")
		viper.AddConfigPath(".")
	}

	// Config can be overwritten by using environment variables prefixed with
	viper.SetEnvPrefix("CONFIG")
	viper.AutomaticEnv()

	// Convert foo.bar into foo_bar so all configuration can be overwritten
	// through environment variable
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		if viper.ConfigFileUsed() == "" {
			log.Fatalf("Unable to find configuration file.")
		}

		log.Fatalf("Failed to load %s: %v", viper.ConfigFileUsed(), err)
	} else {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}
