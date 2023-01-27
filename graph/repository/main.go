package repository

import (
	"context"

	"github.com/jgero/hmc-datastore/graph/model"
)

type Repository interface {
    GetArticles(context.Context) ([]*model.Article, error)
    WriteArticle(context.Context, *model.NewArticle) (*model.Article, error)
}
