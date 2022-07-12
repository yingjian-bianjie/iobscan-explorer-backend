package config

import (
	"bytes"
	"github.com/spf13/viper"
)

type Config struct {
	App     App
	Mongodb MongoDB
}

type App struct {
	Addr            string
	LogLevel        string `mapstructure:"log_level"`
	PrometheusPort  int    `mapstructure:"prometheus_port"`
	ApiKey          string `mapstructure:"api_key"`
	EnableSignature bool   `mapstructure:"enable_signature"`
}

type MongoDB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
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
