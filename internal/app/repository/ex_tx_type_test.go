package repository

import (
	"testing"
)

func Test_exTxTypeRepo_QueryTxTypeList(t *testing.T) {
	list, err := ExTxTypeRepo.QueryTxTypeList()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}
