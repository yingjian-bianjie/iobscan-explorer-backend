package cmd

import (
	"io/ioutil"
	"os"

	app "github.com/bianjieai/iobscan-explorer-backend/internal"
	conf "github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/zk"
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

	testCmd = &cobra.Command{ // test
		Use:   "test",
		Short: "Start iobscan-explorer-backend Server.",
		Run: func(cmd *cobra.Command, args []string) {
			test()
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&localConfig, "CONFIG", "c", defaultLocalConfig, "conf path: /opt/local.toml")
}

func test() {
	config := localConf()
	app.Serve(config)
}

func start() {
	config := onlineConf()
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

func onlineConf() *conf.Config {
	var config *conf.Config

	zkConn, err := zk.NewZkConn()
	if err != nil {
		panic(err)
	}

	configPath := constant.DefaultZkConfigPath
	if v, ok := os.LookupEnv(constant.EnvNameZkConfigPath); ok {
		configPath = v
	}
	data, _, err := zkConn.Get(configPath)
	if err != nil {
		panic(err)
	}
	config, err = conf.ReadConfig(data)
	if err != nil {
		panic(err)
	}
	return config
}
