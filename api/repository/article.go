package repository

import (
	"articlewithgraphql/graph/model"
	"context"

	"github.com/jackc/pgx/v5"
)

func AddArticle(pgx *pgx.Conn, article model.AddArticle) (string, error) {
	_, err := pgx.Exec(context.Background(), `INSERT INTO articles (title, content, image, topic, author) VALUES ($1, $2, $3, $4, $5)`, article.Title, article.Content, article.Image, article.Topic, article.Author)
	if err != nil {
		return "", err
	}
	return "Article Added Successfully", nil
}

func GetMyArticles(pgx *pgx.Conn, author int64) ([]*model.Article, error) {
	articles, err := pgx.Query(context.Background(), `SELECT * FROM articles WHERE author = $1`, author)

	articlesSlice := make([]*model.Article, 0)
	if err != nil {
		return articlesSlice, err
	}
	defer articles.Close()

	var article model.Article
	for articles.Next() {
		if err := articles.Scan(&article.ID, &article.Title, &article.Content, &article.Image, &article.Nooflikes, &article.Noofviews, &article.Topic, &article.Author, &article.Publishedat); err != nil {
			return articlesSlice, err
		}
		articlesSlice = append(articlesSlice, &article)
	}

	return articlesSlice, nil
}

func GetArticleById(pgx *pgx.Conn, id int64) (model.Article, error) {
	row := pgx.QueryRow(context.Background(), `SELECT * FROM articles WHERE id = $1`, id)

	var responseArticle model.Article

	err := row.Scan(&responseArticle.ID, &responseArticle.Title, &responseArticle.Content, &responseArticle.Image, &responseArticle.Nooflikes, &responseArticle.Noofviews, &responseArticle.Topic, &responseArticle.Author, &responseArticle.Publishedat)
	if err != nil {
		return responseArticle, err
	}
	return responseArticle, nil
}

func UpdateArticle(pgx *pgx.Conn, article model.UpdateArticle) (string, error) {
	_, err := pgx.Exec(context.Background(), `UPDATE articles SET title = $1, content = $2, image = $3, topic = $4 WHERE id = $5`, article.Title, article.Content, article.Image, article.Topic, article.ID)
	if err != nil {
		return "", err
	}
	return "Article Updated Successfully.", nil
}

func DeleteArticle(pgx *pgx.Conn, id int64) (string, error) {
	_, err := pgx.Exec(context.Background(), `DELETE FROM articles WHERE id = $1`, id)
	if err != nil {
		return "", err
	}
	return "Article Deleted Successfully.", nil
}
