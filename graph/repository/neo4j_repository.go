package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jgero/hmc-datastore/graph/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jRepo struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jRepo(conn string, user string, password string) Repository {
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "neo4jneo4j", ""),
	)
	if err != nil {
		panic(err)
	}
	return &Neo4jRepo{driver}
}

func (r *Neo4jRepo) Close(ctx context.Context) {
	r.driver.Close(ctx)
}

func (r *Neo4jRepo) GetPosts(ctx context.Context, limit int64, skip int64, keywords []string) ([]*model.Post, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteRead(ctx, session, func(tx neo4j.ManagedTransaction) ([]*model.Post, error) {
		var query string
		if len(keywords) > 0 {
			query = `MATCH (n:Post)--(k:Keyword)
                WITH n, collect(k) as kNodes
                WHERE ALL(kwd IN $keywords WHERE ANY(kn IN kNodes WHERE kn.value = kwd))
                RETURN n ORDER BY n.created DESC SKIP $skip LIMIT $limit`
		} else {
			query = "MATCH (n:Post) RETURN n ORDER BY n.created DESC SKIP $skip LIMIT $limit"
		}
		result, err := tx.Run(ctx, query,
			map[string]any{
				"keywords": keywords,
				"skip":     skip,
				"limit":    limit,
			})
		if err != nil {
			return nil, err
		}
		return extractPostsFromResult(result, ctx)
	})
}

func (r *Neo4jRepo) NewPost(ctx context.Context, a *model.NewPost) (*model.Post, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Post, error) {
		id := uuid.NewString()
		created := time.Now().UnixMilli()
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

func (r *Neo4jRepo) UpdatePost(ctx context.Context, a *model.UpdatePost) (*model.Post, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Post, error) {
		updated := time.Now().UnixMilli()
		query := "MATCH (n:Post { uuid: $uuid }) SET "
		sets := make([]string, 0)
		sets = append(sets, "n.updated = $updated")
		sets = append(sets, "n.updateCount = n.updateCount + 1")
		if a.Title != nil {
			sets = append(sets, "n.title = $title")
		}
		if a.Content != nil {
			sets = append(sets, "n.content = $content")
		}
		query += strings.Join(sets, ", ")
		query += " RETURN n"
		records, err := tx.Run(ctx, query,
			map[string]any{
				"title":   a.Title,
				"content": a.Content,
				"uuid":    a.UUID,
				"updated": updated,
			})
		if err != nil {
			return nil, err
		}
		record, err := records.Single(ctx)
		if err != nil {
			return nil, err
		}
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

func (r *Neo4jRepo) NewPerson(ctx context.Context, p *model.NewPerson) (*model.Person, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Person, error) {
		id := uuid.NewString()
		created := time.Now().UnixMilli()
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

func (r *Neo4jRepo) UpdatePerson(ctx context.Context, a *model.UpdatePerson) (*model.Person, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	return neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (*model.Person, error) {
		updated := time.Now().UnixMilli()
		query := "MATCH (n:Person { uuid: $uuid }) SET "
		sets := make([]string, 0)
		sets = append(sets, "n.updated = $updated")
		sets = append(sets, "n.updateCount = n.updateCount + 1")
		if a.Name != nil {
			sets = append(sets, "n.name = $name")
		}
		query += strings.Join(sets, ", ")
		query += " RETURN n"
		records, err := tx.Run(ctx, query,
			map[string]any{
				"name":    a.Name,
				"uuid":    a.UUID,
				"updated": updated,
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
