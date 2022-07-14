package repository

import (
	"testing"
)

func Test_syncTxRepo_QueryTxMsgsCountByHeight(t *testing.T) {
	height, err := SyncTxRepo.QueryTxMsgsCountByHeight(16481)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(height)
}

func Test_syncTxRepo_QueryTxMsgsIncre(t *testing.T) {
	height, err := SyncTxRepo.QueryTxMsgsIncre(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(height)
}
