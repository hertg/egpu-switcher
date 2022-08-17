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
	Use:          "egpu-switcher",
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

const configPath = "/etc/egpu-switcher"

var verbose bool

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

func Execute() {

	// header := &doc.GenManHeader{
	// 	Title:   "egpu-switcher",
	// 	Section: "1",
	// 	Source:  "egpu-switcher-0.0.1",
	// 	Manual:  "",
	// }

	// err := doc.GenManTree(rootCmd, header, "./docs")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	rootCheck()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCheck() {
	u, err := user.Current()
	if err != nil {
		fmt.Println("unable to get current user. if you run into permission issues, re-try running as root")
		return
	}
	if u.Uid != "0" {
		fmt.Println("please run as root")
		os.Exit(1)
	}
}