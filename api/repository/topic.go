package repository

import (
	"articlewithgraphql/dataloader"
	"articlewithgraphql/graph/model"
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddTopic(pgx *pgxpool.Pool, topic model.AddTopic) (string, error) {
	_, err := pgx.Exec(context.Background(), `INSERT INTO topics (name) VALUES ($1)`, topic.Name)
	if err != nil {
		return "", err
	}
	return "Topic Added Successfully.", nil
}

func GetAllTopics(ctx context.Context, topicIds []string) ([]*model.Topic, error) {
	return dataloader.GetTopics(ctx, topicIds)
}

func GetAllTopics1(pgx *pgxpool.Pool) ([]*model.Topic, error) {
	topics, err := pgx.Query(context.Background(), `SELECT * FROM topics`)
	topicsSlice := make([]*model.Topic, 0)

	if err != nil {
		return topicsSlice, err
	}
	defer topics.Close()

	for topics.Next() {
		var topic model.Topic
		if err := topics.Scan(&topic.ID, &topic.Name); err != nil {
			return topicsSlice, err
		}
		topicsSlice = append(topicsSlice, &topic)
	}

	return topicsSlice, nil
}

func GetArticlesByTopicId1(ctx context.Context, topicId *int64) ([]*model.Article, error) {

	allArticles, err := dataloader.GetArticles(ctx, []string{strconv.FormatInt(*topicId, 10)})

	fmt.Println(err)

	return allArticles, err

}

func GetArticlesByTopicId(ctx context.Context, pgx *pgxpool.Pool, topicId int64) ([]*model.Article, error) {

	rows, err := pgx.Query(ctx, `select id, title, content, topic, author from articles where topic = $1`, topicId)
	// topics, err := pgx.Query(context.Background(), `SELECT * FROM topics`)
	articlesSlice := make([]*model.Article, 0)

	if err != nil {
		return articlesSlice, err
	}
	defer rows.Close()
	fmt.Println(topicId)

	for rows.Next() {
		fmt.Println(topicId)
		var article model.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Topic, &article.Author); err != nil {
			return articlesSlice, err
		}
		articlesSlice = append(articlesSlice, &article)
	}

	return articlesSlice, nil
}

func GetArticlesByTopicId2(ctx context.Context, pgx *pgxpool.Pool, topicId int64) ([]*model.Article, error) {

	rows, err := pgx.Query(ctx, `select id, title, content, topic, author from articles where topic = $1`, topicId)

	if err != nil {
		return nil, err
	}

	var articlesSlice []*model.Article
	for rows.Next() {
		var article model.Article
		rows.Scan(&article.ID, &article.Title, &article.Content, &article.Topic, &article.Author)
		articlesSlice = append(articlesSlice, &article)
	}
	return articlesSlice, nil
}

func UpdateTopic(pgx *pgxpool.Pool, topic model.UpdateTopic) (string, error) {
	_, err := pgx.Exec(context.Background(), `UPDATE topics SET name = $1 WHERE id = $2`, topic.Name, topic.ID)
	if err != nil {
		return "", err
	}
	return "Topic Updated Successfully.", nil
}

func DeleteTopic(pgx *pgxpool.Pool, id int64) (string, error) {
	_, err := pgx.Exec(context.Background(), `DELETE FROM topics WHERE id = $1`, id)
	if err != nil {
		return "", err
	}
	return "Topic Deleted Successfully.", nil
}
