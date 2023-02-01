package repository

import (
	"context"
	"fmt"
	"time"

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

func (r *Neo4jRepo) GetPosts(ctx context.Context) ([]*model.Post, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteRead(ctx, session, func(tx neo4j.ManagedTransaction) ([]*model.Post, error) {
		result, err := tx.Run(ctx, "MATCH (n:Post) RETURN n", map[string]any{})
		if err != nil {
			return nil, err
		}
		return extractPostsFromResult(result, ctx)
	})
}

func (r *Neo4jRepo) WritePost(ctx context.Context, a *model.NewPost) (*model.Post, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Post, error) {
		id := uuid.NewString()
		created := time.Now().Unix()
		records, err := tx.Run(ctx, `
            MATCH (p:Person { uuid: $writerUuid })
            CREATE (n:Post { title: $title, content: $content, uuid: $uuid, created: $created, updated: $created, updateCount: 0 })<-[:writer]-(p)
            WITH n
            FOREACH (kwd in $keywords |
              MERGE (k:Keyword {value:kwd})
              MERGE (n)-[:relates_to]->(k)
            )
            RETURN n`,
			map[string]any{
				"title":      a.Title,
				"content":    a.Content,
				"uuid":       id,
				"created":    created,
				"writerUuid": a.WriterUUID,
				"keywords":   a.Keywords,
			})
		if err != nil {
			return nil, err
		}
		record, err := records.Single(ctx)
		return extractPostFromRecord(record, ctx)
	})
}

func (r *Neo4jRepo) GetWriter(ctx context.Context, a *model.Post) (*model.Person, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteRead(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Person, error) {
		result, err := tx.Run(ctx, `MATCH (:Post { uuid: $uuid })<-[:writer]-(n:Person) RETURN n`,
			map[string]any{
				"uuid": a.UUID,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return extractPersonFromRecord(record, ctx)
	})
}

func (r *Neo4jRepo) WritePerson(ctx context.Context, p *model.NewPerson) (*model.Person, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Person, error) {
		id := uuid.NewString()
		created := time.Now().Unix()
		records, err := tx.Run(ctx, `
            CREATE (n:Person { uuid: $uuid, name: $name, created: $created, updated: $created, updateCount: 0 })
            WITH n
            FOREACH (kwd in $keywords |
                MERGE (k:Keyword {value:kwd})
                MERGE (n)-[:relates_to]->(k)
            )
            RETURN n`,
			map[string]any{
				"uuid":     id,
				"name":     p.Name,
				"keywords": p.Keywords,
				"created":  created,
			})
		if err != nil {
			return nil, err
		}
		record, err := records.Single(ctx)
		if err != nil {
			return nil, err
		}
		return extractPersonFromRecord(record, ctx)
	})
}

func (r *Neo4jRepo) GetKeywordsForUuid(ctx context.Context, k string) ([]*model.Keyword, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteRead(ctx, session, func(tx neo4j.ManagedTransaction) ([]*model.Keyword, error) {
		result, err := tx.Run(ctx, `
            MATCH ({uuid: $uuid})-[:relates_to]->(k:Keyword)
            OPTIONAL MATCH (k)<-[r:relates_to]-()
            RETURN k.value, count(r) AS usages ORDER BY usages DESC`,
			map[string]any{
				"uuid": k,
			})
		if err != nil {
			return nil, err
		}
		return extractKeywordsFromResult(result, ctx)
	})
}

func (r *Neo4jRepo) GetKeywords(ctx context.Context) ([]*model.Keyword, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteRead(ctx, session, func(tx neo4j.ManagedTransaction) ([]*model.Keyword, error) {
		result, err := tx.Run(ctx, `
            MATCH (k:Keyword) OPTIONAL MATCH (k)<-[r:relates_to]-()
            RETURN k.value, count(r) AS usages ORDER BY usages DESC`,
			map[string]any{})
		if err != nil {
			return nil, err
		}
		return extractKeywordsFromResult(result, ctx)
	})
}

func (r *Neo4jRepo) WriteKeywords(ctx context.Context, s *model.SetKeywords) ([]model.KeywordLink, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) ([]model.KeywordLink, error) {
		query := `
            MATCH (n) WHERE n.uuid in $uuids
            WITH n, $keywords as kwds
            FOREACH (kwd in kwds |
              MERGE (k:Keyword {value:kwd})
              MERGE (n)-[:relates_to]->(k)
            )`
		if s.Exclusive {
			query += `
                WITH n, kwds
                MATCH (n)-[r:relates_to]-(k:Keyword) 
                WHERE NOT k.value IN kwds
                DELETE r`
		}
		query += `
            RETURN n`
		result, err := tx.Run(ctx, query,
			map[string]any{
				"uuids":    s.UUIDs,
				"keywords": s.Keywords,
			})
		if err != nil {
			return nil, err
		}
		modifiedItems := make([]model.KeywordLink, 0)
		for result.Next(ctx) {
			record := result.Record()
			rawNode, found := record.Get("n")
			if !found {
				return nil, fmt.Errorf("could not find column")
			}
			actualNode := rawNode.(neo4j.Node)
			for _, v := range actualNode.Labels {
				switch v {
				case "Person":
					p, err := extractPersonFromRecord(record, ctx)
					if err != nil {
						return nil, err
					}
					modifiedItems = append(modifiedItems, p)
				case "Post":
					p, err := extractPostFromRecord(record, ctx)
					if err != nil {
						return nil, err
					}
					modifiedItems = append(modifiedItems, p)
				}
			}
		}
		return modifiedItems, nil
	})
}

func extractPersonFromRecord(record *neo4j.Record, ctx context.Context) (*model.Person, error) {
	rawNode, found := record.Get("n")
	if !found {
		return nil, fmt.Errorf("could not find column")
	}
	personNode := rawNode.(neo4j.Node)
	name, err := neo4j.GetProperty[string](personNode, "name")
	if err != nil {
		return nil, err
	}
	uuid, err := neo4j.GetProperty[string](personNode, "uuid")
	if err != nil {
		return nil, err
	}
	created, err := neo4j.GetProperty[int64](personNode, "created")
	if err != nil {
		return nil, err
	}
	updated, err := neo4j.GetProperty[int64](personNode, "updated")
	if err != nil {
		return nil, err
	}
	updateCount, err := neo4j.GetProperty[int64](personNode, "updateCount")
	if err != nil {
		return nil, err
	}
	return &model.Person{Name: name, UUID: uuid, Created: created, Updated: updated, UpdateCount: updateCount}, nil
}

func extractKeywordsFromResult(result neo4j.ResultWithContext, ctx context.Context) ([]*model.Keyword, error) {
	keywords := make([]*model.Keyword, 0)
	for result.Next(ctx) {
		record := result.Record()
		rawKeywordValue, found := record.Get("k.value")
		if !found {
			return nil, fmt.Errorf("could not find column")
		}
		keyword := rawKeywordValue.(string)
		rawUsages, found := record.Get("usages")
		if !found {
			return nil, fmt.Errorf("could not find column")
		}
		usages := rawUsages.(int64)
		keywords = append(keywords, &model.Keyword{Value: keyword, Usages: usages})
	}
	return keywords, nil
}

func extractPostFromRecord(record *neo4j.Record, ctx context.Context) (*model.Post, error) {
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
	created, err := neo4j.GetProperty[int64](articleNode, "created")
	if err != nil {
		return nil, err
	}
	updated, err := neo4j.GetProperty[int64](articleNode, "updated")
	if err != nil {
		return nil, err
	}
	updateCount, err := neo4j.GetProperty[int64](articleNode, "updateCount")
	if err != nil {
		return nil, err
	}
	return &model.Post{Title: title, Content: content, UUID: uuid, Created: created, Updated: updated, UpdateCount: updateCount}, nil
}

func extractPostsFromResult(result neo4j.ResultWithContext, ctx context.Context) ([]*model.Post, error) {
	articles := make([]*model.Post, 0)
	for result.Next(ctx) {
		record := result.Record()
		a, err := extractPostFromRecord(record, ctx)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, nil
}
