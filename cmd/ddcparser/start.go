package main

import (
	"os"

	"github.com/bianjieai/iobscan-explorer-backend/internal/common/configs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultCLIHome = os.ExpandEnv("$HOME/.ddc-parser")
	flagHome       = "home"
)

// StartCmd return the start command
func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Example: "ddcparser start",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := readConfig(cmd)
			if err != nil {
				panic(err)
			}
			start(cfg)
		},
	}
	cmd.Flags().String(flagHome, defaultCLIHome, "ddcparser server config path")
	return cmd
}

func readConfig(cmd *cobra.Command) (*configs.Config, error) {
	rootViper := viper.New()
	_ = rootViper.BindPFlags(cmd.Flags())
	// Find home directory.
	rootViper.AddConfigPath(rootViper.GetString(flagHome))
	rootViper.SetConfigName("config")
	rootViper.SetConfigType("toml")

	// Find and read the config file
	if err := rootViper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return nil, err
	}

	var config configs.Config
	if err := rootViper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
