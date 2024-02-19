package repository

import (
	"articlewithgraphql/api/validation"
	"articlewithgraphql/error"
	"articlewithgraphql/graph/model"
	"articlewithgraphql/utils"
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUser(pgx *pgxpool.Pool, user model.RegisterUser) (string, error) {
	isEmail := validation.EmailValidation(user.Email)
	if !isEmail {
		return "", errorhandling.EmailvalidationError
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashedPassword

	_, err = pgx.Exec(context.Background(), `INSERT INTO users (name, bio, email, password) VALUES ($1, $2, $3, $4)`, user.Name, user.Bio, user.Email, user.Password)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			return "", errorhandling.DuplicateEmailFound
		}
		return "", errorhandling.RegistrationFailedError
	}
	return "Registration Done Successfully.", nil

}

func LoginUser(pgx *pgxpool.Pool, user model.LoginUser) (model.User, error) {
	var dbUser model.User
	row := pgx.QueryRow(context.Background(), `SELECT id, name, bio, email, password, image, isadmin FROM users WHERE email = $1`, user.Email)
	err := row.Scan(&dbUser.ID, &dbUser.Name, &dbUser.Bio, &dbUser.Email, &dbUser.Password, &dbUser.Image, &dbUser.Isadmin)

	if err == sql.ErrNoRows {
		return dbUser, errorhandling.NoUserFound
	}

	// passwordMatched := utils.VerifyPassword(user.Password, dbUser.Password)
	// if !passwordMatched {
	// 	return dbUser, errorhandling.PasswordNotMatch
	// }

	accesstoken, err := utils.CreateAccessToken(time.Now().Add(time.Hour*24*7), *dbUser.ID, *dbUser.Isadmin)
	if err != nil {
		return dbUser, err
	}

	dbUser.AccessToken = &model.AccessToken{Token: accesstoken}

	return dbUser, nil
}

//scalar
//error
//go loading
