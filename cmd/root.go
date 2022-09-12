package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:               "egpu-switcher",
	Short:             "Distribution agnostic eGPU script that works with NVIDIA and AMD cards.",
	SilenceUsage:      true,
	DisableAutoGenTag: true,

	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

const configPath = "/etc/egpu-switcher"

var verbose bool
var isRoot bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func initConfig() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	// map environment variables with underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// defaults
	viper.SetDefault("detection.retries", 6)
	viper.SetDefault("detection.interval", 500)

	u, err := user.Current()
	if err != nil {
		logger.Warn("unable to get current user. if you run into permission issues, retry running as root")
		isRoot = true // just assume :)
	} else if u.Uid == "0" {
		isRoot = true
	}

	// only read in / write default config when user is root
	// simply ignore and don't load config if non-root
	if isRoot {
		err := viper.ReadInConfig()
		if err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				if verbose {
					logger.Debug("no configuration file found, creating a new one at %s\n", configPath)
				}
				err = os.MkdirAll(configPath, 0744)
				cobra.CheckErr(err)
				err = viper.SafeWriteConfig()
				cobra.CheckErr(err)
			default:
				fmt.Println("unable to load config:", err)
				os.Exit(1)
			}
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
