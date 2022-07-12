package cmd

import (
	"io/ioutil"

	app "github.com/bianjieai/iobscan-explorer-backend/internal"
	conf "github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
	"github.com/spf13/cobra"
)

const (
	defaultLocalConfig = "/opt/iobscan-explorer-backend/cfg.toml"
)

var (
	localConfig string
	startCmd    = &cobra.Command{
		Use:   "start",
		Short: "Start iobscan-explorer-backend Server.",
		Run: func(cmd *cobra.Command, args []string) {
			start()
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&localConfig, "CONFIG", "c", defaultLocalConfig, "conf path: /opt/local.toml")
}

func start() {
	config := localConf()
	app.Serve(config)
}

func localConf() *conf.Config {
	if localConfig == "" {
		localConfig = defaultLocalConfig
	}
	data, err := ioutil.ReadFile(localConfig)
	if err != nil {
		panic(err)
	}
	config, err := conf.ReadConfig(data)
	if err != nil {
		panic(err)
	}
	return config
}
