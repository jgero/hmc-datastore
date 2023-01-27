package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jgero/hmc-datastore/graph/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jRepo struct {
	driver neo4j.DriverWithContext
}

var repo *Neo4jRepo

func GetNeo4jRepo() Repository {
	if repo == nil {
		driver, err := neo4j.NewDriverWithContext(
			"neo4j://localhost:7687",
			neo4j.BasicAuth("neo4j", "neo4jneo4j", ""),
		)
		if err != nil {
			panic(err)
		}
        // TODO: fix closing driver
		// ctx := context.Background()
		// close driver when background shuts down
		// defer driver.Close(ctx)
		repo = &Neo4jRepo{driver}
	}
	return repo
}

func (r *Neo4jRepo) GetArticles(ctx context.Context) ([]*model.Article, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteRead(ctx, session, func(tx neo4j.ManagedTransaction) ([]*model.Article, error) {
		result, err := tx.Run(ctx, "MATCH (n:Article) RETURN n", map[string]any{})
		if err != nil {
			return nil, err
		}
		articles := make([]*model.Article, 0)
		for result.Next(ctx) {
			record := result.Record()
			rawNode, found := record.Get("n")
			if !found {
				return nil, fmt.Errorf("could not find column")
			}
			articleNode := rawNode.(neo4j.Node)
			title, err := neo4j.GetProperty[string](articleNode, "title")
			if err != nil {
				return nil, err
			}
			content, err := neo4j.GetProperty[string](articleNode, "content")
			if err != nil {
				return nil, err
			}
			articles = append(articles, &model.Article{Title: title, Content: content, Writer: nil})
		}
		return articles, nil
	})
}

func (r *Neo4jRepo) WriteArticle(ctx context.Context, a *model.NewArticle) (*model.Article, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Article, error) {
		id := uuid.NewString()
		records, err := tx.Run(ctx, "CREATE (n:Article { title: $title, content: $content, uuid: $uuid }) RETURN n", map[string]any{
			"title":   a.Title,
			"content": a.Content,
			"uuid":    id,
		})
		// In face of driver native errors, make sure to return them directly.
		// Depending on the error, the driver may try to execute the function again.
		if err != nil {
			return nil, err
		}
		record, err := records.Single(ctx)
		if err != nil {
			return nil, err
		}
		rawNode, found := record.Get("n")
		if !found {
			return nil, fmt.Errorf("could not find column")
		}
		articleNode := rawNode.(neo4j.Node)
		title, err := neo4j.GetProperty[string](articleNode, "title")
		if err != nil {
			return nil, err
		}
		content, err := neo4j.GetProperty[string](articleNode, "content")
		if err != nil {
			return nil, err
		}
		return &model.Article{Title: title, Content: content, Writer: nil}, nil
	})
}
