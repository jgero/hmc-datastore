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
			uuid, err := neo4j.GetProperty[string](articleNode, "uuid")
			if err != nil {
				return nil, err
			}
			articles = append(articles, &model.Article{Title: title, Content: content, Uuid: uuid})
		}
		return articles, nil
	})
}

func (r *Neo4jRepo) WriteArticle(ctx context.Context, a *model.NewArticle) (*model.Article, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Article, error) {
		id := uuid.NewString()
		records, err := tx.Run(ctx, `
            MATCH (p:Person { uuid: $writerUuid })
            CREATE (n:Article { title: $title, content: $content, uuid: $uuid })<-[:writer]-(p)
            RETURN n`,
			map[string]any{
				"title":      a.Title,
				"content":    a.Content,
				"uuid":       id,
				"writerUuid": a.WriterUUID,
			})
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
		uuid, err := neo4j.GetProperty[string](articleNode, "uuid")
		if err != nil {
			return nil, err
		}
		return &model.Article{Title: title, Content: content, Uuid: uuid}, nil
	})
}

func (r *Neo4jRepo) GetWriter(ctx context.Context, a *model.Article) (*model.Person, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteRead(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Person, error) {
		result, err := tx.Run(ctx, `MATCH (:Article { uuid: $uuid })<-[:writer]-(p:Person) RETURN p`,
			map[string]any{
				"uuid": a.Uuid,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		rawNode, found := record.Get("p")
		if !found {
			return nil, fmt.Errorf("could not find column")
		}
		writerNode := rawNode.(neo4j.Node)
		name, err := neo4j.GetProperty[string](writerNode, "name")
		if err != nil {
			return nil, err
		}
		uuid, err := neo4j.GetProperty[string](writerNode, "uuid")
		if err != nil {
			return nil, err
		}
		return &model.Person{Name: name, UUID: uuid}, nil
	})
}

func (r *Neo4jRepo) WritePerson(ctx context.Context, p *model.NewPerson) (*model.Person, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Person, error) {
		id := uuid.NewString()
		records, err := tx.Run(ctx, `CREATE (p:Person { uuid: $uuid, name: $name }) RETURN p`,
			map[string]any{
				"uuid": id,
				"name": p.Name,
			})
		if err != nil {
			return nil, err
		}
		record, err := records.Single(ctx)
		if err != nil {
			return nil, err
		}
		rawNode, found := record.Get("p")
		if !found {
			return nil, fmt.Errorf("could not find column")
		}
		personNode := rawNode.(neo4j.Node)
		uuid, err := neo4j.GetProperty[string](personNode, "uuid")
		if err != nil {
			return nil, err
		}
		name, err := neo4j.GetProperty[string](personNode, "name")
		if err != nil {
			return nil, err
		}
		return &model.Person{Name: name, UUID: uuid}, nil
	})
}

func (r *Neo4jRepo) GetKeywords(ctx context.Context, k string) ([]string, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteRead(ctx, session, func(tx neo4j.ManagedTransaction) ([]string, error) {
		result, err := tx.Run(ctx, `MATCH (:Article {uuid: $uuid})-[:relates_to]->(k:Keyword) RETURN k.value`,
			map[string]any{
				"uuid": k,
			})
		if err != nil {
			return nil, err
		}
		keywords := make([]string, 0)
		for result.Next(ctx) {
			record := result.Record()
			rawNode, found := record.Get("k.value")
			if !found {
				return nil, fmt.Errorf("could not find column")
			}
			keyword := rawNode.(string)
			keywords = append(keywords, keyword)
		}
		return keywords, nil
	})
}
