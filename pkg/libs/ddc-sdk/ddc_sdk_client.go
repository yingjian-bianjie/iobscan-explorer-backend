package ddc_sdk

import (
	"github.com/bianjieai/ddc-sdk-go/ddc-sdk-platform-go/app"
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/configs"
)

var (
	_client *app.DDCSdkClient
)

func Client() *app.DDCSdkClient {
	if _client != nil {
		return _client
	}
	panic("InitDDCSDKClient not work")
}

func InitDDCSDKClient(conf configs.DdcClientConf) {
	clientBuilder := app.DDCSdkClientBuilder{}
	client := clientBuilder.
		SetSignEventListener(new(SignListener)). //注册实现了签名接口的结构体
		//SetGasPrice(conf.GasPrice).              //建议设置gasPrice>=1e10
		SetGatewayURL(conf.GatewayURL).
		SetAuthorityAddress(conf.AuthorityAddress).
		SetChargeAddress(conf.ChargeAddress).
		SetDDC721Address(conf.DDC721Address).
		SetDDC1155Address(conf.DDC1155Address)
	if conf.GatewayAPIKey != "" && conf.GatewayAPIValue != "" {
		client = client.SetGatewayAPIKey(conf.GatewayAPIKey).SetGatewayAPIValue(conf.GatewayAPIValue)
	}
	if conf.LogFilepath != "" {
		client = client.RegisterLog(conf.LogFilepath)
	}
	_client = client.Build()
}
