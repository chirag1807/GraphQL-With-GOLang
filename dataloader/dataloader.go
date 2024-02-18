package dataloader

import (
	"articlewithgraphql/graph/model"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vikstrous/dataloadgen"
)

type ctxKeyType string

const (
	loadersKey ctxKeyType = "dataloaders"
)

type DbReader struct {
	pgx *pgx.Conn
}

type Loaders struct {
	TopicArticlesLoader *dataloadgen.Loader[string, *model.Article]
	TopicsLoader        *dataloadgen.Loader[string, *model.Topic]
}

func NewLoaders(pgx *pgx.Conn) *Loaders {
	dr := &DbReader{pgx: pgx}
	return &Loaders{
		TopicArticlesLoader: dataloadgen.NewLoader(dr.GetTopicArticles, dataloadgen.WithWait(time.Millisecond)),
		TopicsLoader:        dataloadgen.NewLoader(dr.GetTopics, dataloadgen.WithWait(time.Millisecond)),
	}
}

func DataLoaderMiddleware(pgx *pgx.Conn, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(pgx)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

func (c DbReader) GetTopicArticles(ctx context.Context, topicIds []string) ([]*model.Article, []error) {
	var idInterfaces []interface{}
	for _, id := range topicIds {
		idInterfaces = append(idInterfaces, id)
	}

	placeholders := make([]string, len(topicIds))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := `select id, title, content from articles where topic in (` + strings.Join(placeholders, ",") + `)`

	rows, err := c.pgx.Query(ctx, query, idInterfaces...)

	if err != nil {
		return nil, []error{err}
	}

	defer rows.Close()

	// articles := make([]*model.Article, 0, len(topicIds))
	var articles []*model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Content); err != nil {
			// errs = append(errs, err)
			// continue
			return nil, []error{err}
		}
		articles = append(articles, &article)
	}

	if err := rows.Err(); err != nil {
		return nil, []error{err}
	}

	return articles, nil
}

func GetArticle(ctx context.Context, topicId string) (*model.Article, error) {
	loaders := For(ctx)
	return loaders.TopicArticlesLoader.Load(ctx, topicId)
}

func GetArticles(ctx context.Context, topicIds []string) ([]*model.Article, error) {
	loaders := For(ctx)
	return loaders.TopicArticlesLoader.LoadAll(ctx, topicIds)
}

func (c DbReader) GetTopics(ctx context.Context, topicIds []string) ([]*model.Topic, []error) {
	var idInterfaces []interface{}
	for _, id := range topicIds {
		idInterfaces = append(idInterfaces, id)
	}

	placeholders := make([]string, len(topicIds))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := `select id, name from topics where id in (` + strings.Join(placeholders, ",") + `)`
	rows, err := c.pgx.Query(ctx, query, idInterfaces...)

	// query := `select id, name from topics`
	// rows, err := c.pgx.Query(ctx, query)

	if err != nil {
		return nil, []error{err}
	}

	var topics []*model.Topic
	errs := make([]error, 0, len(topicIds))
	for rows.Next() {
		// var topic model.Topic
		topic := &model.Topic{}
		if err := rows.Scan(&topic.ID, &topic.Name); err != nil {
			fmt.Println(err)
			errs = append(errs, err)
			continue
		}
		topics = append(topics, topic)
	}
	if len(errs) == 0 {
		return topics, nil
	}
	return topics, errs
}

func GetTopic(ctx context.Context, topicId string) (*model.Topic, error) {
	loaders := For(ctx)
	return loaders.TopicsLoader.Load(ctx, topicId)
}

func GetTopics(ctx context.Context, topicIds []string) ([]*model.Topic, error) {
	loaders := For(ctx)
	return loaders.TopicsLoader.LoadAll(ctx, topicIds)
}
