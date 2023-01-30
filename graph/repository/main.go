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
    GetKeywords(context.Context, string) ([]string, error)
    WriteKeywords(context.Context, *model.SetKeywords) ([]string, error)
}
