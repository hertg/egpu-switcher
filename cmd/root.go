package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "egpu-switcher",
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

const configPath = "/etc/egpu-switcher"

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

func initConfig() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	// default values
	viper.SetDefault("verbose", false)

	// bind cobra flags to viper config
	viper.BindPFlags(rootCmd.Flags())

	// map environment variables with underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// todo: create config file
			fmt.Println("todo: create config file")
			os.MkdirAll(configPath, 0744)
			err = viper.SafeWriteConfig()
			cobra.CheckErr(err)
			os.Exit(1)
		default:
			fmt.Println("unable to load config:", err)
			os.Exit(1)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
