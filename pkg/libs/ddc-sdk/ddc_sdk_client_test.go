package ddc_sdk

import (
	"fmt"
	"testing"

	"github.com/bianjieai/iobscan-explorer-backend/internal/common/configs"
)

func TestMain(m *testing.M) {
	InitDDCSDKClient(configs.DdcClientConf{
		//GasPrice:         1,
		GatewayURL:       "http://192.168.150.42:8545",
		AuthorityAddress: "0xBcE9AA1924D7197C9C945e43638Bf589f91bcB71",
		ChargeAddress:    "0xF41b6185bFB22E2EFC5fB8395Fa3B952951E2d0b",
		DDC721Address:    "0x74b6114d011891Ac21FD1d586bc7F3407c63c216",
		DDC1155Address:   "0x9f7388e114DfDFAbAF8e4b881894E4C7e1b52C17",
		//GatewayURL:       "http://47.100.192.234:8545",
		//AuthorityAddress: "0xa7FC5B0F4A0085c5Ce689b919a866675Ce37B66b",
		//ChargeAddress:    "0xF41b6185bFB22E2EFC5fB8395Fa3B952951E2d0b",
		//DDC721Address:    "0x3B09b7A00271C5d9AE84593850dE3A526b8BF96e",
		//DDC1155Address:   "0xe5d3b9E7D16E03A4A1060c72b5D1cb7806DD9070",
		LogFilepath: "./log.log",
	})
	m.Run()
}

func TestGetDDCId(t *testing.T) {

	fmt.Println(_client.GetDDC721Service().DDCIdByHash("0xbcd980cf6d3e0493ebaaeb72054421735c1f3bb04bcd621c44854bf25ee88745"))
}

func TestMethod(t *testing.T) {
	name, _ := _client.GetDDC721Service().Name()
	owner, _ := _client.GetDDC721Service().OwnerOf(518)
	symbol, _ := _client.GetDDC721Service().Symbol()
	uri, _ := _client.GetDDC721Service().DdcURI(518)
	fmt.Println(name, owner, symbol, uri)
	address, _ := _client.GetDDC721Service().HexToBech32(owner)
	fmt.Println("iaaxxx:", address)
	oXaddr, _ := _client.GetDDC721Service().Bech32ToHex(address)
	fmt.Println("0xXXX:", oXaddr)
}

func TestOwnerOf(t *testing.T) {
	owner, err := _client.GetDDC721Service().OwnerOf(1835)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(owner)
}
