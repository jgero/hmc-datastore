package repository

import (
	"context"

	"github.com/jgero/hmc-datastore/graph/model"
)

type Repository interface {
    Close(context.Context)
    GetPosts(context.Context, int64, int64, []string) ([]*model.Post, error)
    NewPost(context.Context, *model.NewPost) (*model.Post, error)
    UpdatePost(context.Context, *model.UpdatePost) (*model.Post, error)
    NewPerson(context.Context, *model.NewPerson) (*model.Person, error)
    UpdatePerson(context.Context, *model.UpdatePerson) (*model.Person, error)
    GetWriter(context.Context, *model.Post) (*model.Person, error)
    GetKeywordsForUuid(context.Context, string) ([]*model.Keyword, error)
    GetKeywords(context.Context) ([]*model.Keyword, error)
    WriteKeywords(context.Context, *model.SetKeywords) ([]model.KeywordLink, error)
}
