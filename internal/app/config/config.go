package config

import (
	"bytes"

	"github.com/spf13/viper"
)

type Config struct {
	App        App
	Mongodb    MongoDB
	Lcd        Lcd
	BlockChain BlockChain
	Task       Task
}

type App struct {
	Addr            string
	LogLevel        string `mapstructure:"log_level"`
	PrometheusPort  int    `mapstructure:"prometheus_port"`
	ApiKey          string `mapstructure:"api_key"`
	EnableSignature bool   `mapstructure:"enable_signature"`
}

type MongoDB struct {
	Url      string `mapstructure:"url"`
	Database string `mapstructure:"database"`
}

type Lcd struct {
	Backend          string `mapstructure:"backend"`
	BondTokensUrl    string `mapstructure:"bond_tokens_url"`
	TotalSupplyUrl   string `mapstructure:"total_supply_url"`
	CommunityPoolUrl string `mapstructure:"community_pool_url"`
}

type BlockChain struct {
	RpcAddr  string `mapstructure:"rpc_addr"`
	GrpcAddr string `mapstructure:"grpc_addr"`
	ChainId  string `mapstructure:"chain_id"`
}

type Task struct {
	CronJobs string `mapstructure:"cron_jobs"`
}

func ReadConfig(data []byte) (*Config, error) {
	v := viper.New()
	v.SetConfigType("toml")
	reader := bytes.NewReader(data)
	err := v.ReadConfig(reader)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
