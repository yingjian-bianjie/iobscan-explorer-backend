package repository

import (
	"testing"
)

func Test_syncBlockRepo_QueryLatestBlockHeight(t *testing.T) {
	height, err := SyncBlockRepo.QueryLatestBlockHeight(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(height.Height)
}
