package db

import (
	"github.com/gocql/gocql"
	"github.com/rmortale/bookstore_oauth-api/src/clients/cassandra"
	"github.com/rmortale/bookstore_oauth-api/src/domain/access_token"
	"github.com/rmortale/bookstore_utils-go/rest_errors"
)

const (
	queryGetAccessToken    = "select access_token, user_id, client_id, expires from access_tokens where access_token=?;"
	queryCreateAccessToken = "insert into access_tokens(access_token, user_id, client_id, expires) values (?,?,?,?);"
	queryUpdateExpires     = "update access_tokens set expires=? where access_token=?;"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessToken) rest_errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, rest_errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken, &result.UserId, &result.ClientId, &result.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError("no access token found with given id")
		}
		return nil, rest_errors.NewInternalServerError(err.Error(), err)
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken, at.UserId, at.ClientId, at.Expires,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error(), err)
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires, at.AccessToken,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error(), err)
	}
	return nil
}
