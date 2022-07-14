package repository

import (
	"context"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/model"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionNameExTokens = "ex_tokens"
)

func NewExTokensRepo(cli *qmgo.Client, database string) IExTokensRepo {
	return &exTokensRepo{coll: cli.Database(database).Collection(CollectionNameExTokens)}
}

type IExTokensRepo interface {
	QueryMainToken() (*model.Tokens, error)
}

type exTokensRepo struct {
	coll *qmgo.Collection
}

func (repo *exTokensRepo) QueryMainToken() (*model.Tokens, error) {
	var token model.Tokens
	q := bson.M{
		"is_main_token": true,
	}
	err := repo.coll.Find(context.Background(), q).One(&token)
	return &token, err
}
