package configs

type (
	Config struct {
		DataBaseConf DataBaseConf  `mapstructure:"database"`
		DdcClient    DdcClientConf `mapstructure:"ddc_client"`
		Server       Server        `mapstructure:"server"`
	}

	DataBaseConf struct {
		Addrs    string `mapstructure:"addrs"`
		User     string `mapstructure:"user"`
		Passwd   string `mapstructure:"passwd" json:"-"`
		Database string `mapstructure:"database"`
	}
	Server struct {
		Prometheus        string `mapstructure:"prometheus_port"`
		IncreHeight       int64  `mapstructure:"incre_height"`
		InsertBatchLimit  int    `mapstructure:"insert_batch_limit"`
		MaxOperateTxCount int    `mapstructure:"max_operate_tx_count"`
		GasPrice          string `mapstructure:"gas_price"`
	}
	DdcClientConf struct {
		GatewayURL       string `mapstructure:"gateway_url"`
		GatewayAPIKey    string `mapstructure:"gateway_api_key"`
		GatewayAPIValue  string `mapstructure:"gateway_api_value"`
		AuthorityAddress string `mapstructure:"authority_address"`
		ChargeAddress    string `mapstructure:"charge_address"`
		DDC721Address    string `mapstructure:"ddc_721_address"`
		DDC1155Address   string `mapstructure:"ddc_1155_address"`
		LogFilepath      string `mapstructure:"log_filepath"`
	}
)
