package repository

import (
	"context"

	"github.com/jgero/hmc-datastore/graph/model"
)

type Repository interface {
    GetArticles(context.Context) ([]*model.Article, error)
    WriteArticle(context.Context, *model.NewArticle) (*model.Article, error)
    WritePerson(context.Context, *model.NewPerson) (*model.Person, error)
    GetWriter(context.Context, *model.Article) (*model.Person, error)
    GetKeywords(context.Context, string) ([]string, error)
}
