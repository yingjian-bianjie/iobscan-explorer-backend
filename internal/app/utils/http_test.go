package utils

import (
	"testing"
)

func TestHttpGetJson(t *testing.T) {
	url := `http://serviceiris.upticknft.com/metadata/uptickde3073f7fd19c8fa6f9806a3be505c25/uptickc460cc3abc7a4237afd7ec09376372d4.json`
	//url := `http://ipfs.uptick.world:8082/ipfs/QmZwdVJ6jM7ouCELQkytTr1DxTczNPhys4jmphkoca3qGP`
	json, contentType, err := HttpGetJsonOrContentType(url)
	if err != nil {
		t.Fatal(err)
	}
	if contentType == "" {
		t.Log(string(json))
	} else {
		t.Log(contentType)
	}
}
