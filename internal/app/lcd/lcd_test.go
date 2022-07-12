package lcd

import (
	"testing"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/utils"
)

func TestGetBondedTokens(t *testing.T) {
	tokens, err := GetBondedTokens()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.MustMarshalJsonToStr(tokens))
}

func TestGetTotalSupply(t *testing.T) {
	supply, err := GetTotalSupply()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.MustMarshalJsonToStr(supply))
}
