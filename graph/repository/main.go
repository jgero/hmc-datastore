package repository

import (
	"context"

	"github.com/jgero/hmc-datastore/graph/model"
)

type Repository interface {
    GetPosts(context.Context) ([]*model.Post, error)
    WritePost(context.Context, *model.NewPost) (*model.Post, error)
    WritePerson(context.Context, *model.NewPerson) (*model.Person, error)
    GetWriter(context.Context, *model.Post) (*model.Person, error)
    GetKeywordsForUuid(context.Context, string) ([]*model.Keyword, error)
    GetKeywords(context.Context) ([]*model.Keyword, error)
    WriteKeywords(context.Context, *model.SetKeywords) ([]model.KeywordLink, error)
}
