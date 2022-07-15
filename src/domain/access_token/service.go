package access_token

import (
	"github.com/rmortale/bookstore_utils-go/rest_errors"
	"strings"
)

type Repository interface {
	GetById(string) (*AccessToken, rest_errors.RestErr)
	Create(AccessToken) rest_errors.RestErr
	UpdateExpirationTime(AccessToken) rest_errors.RestErr
}

type Service interface {
	GetById(string) (*AccessToken, rest_errors.RestErr)
	Create(AccessToken) rest_errors.RestErr
	UpdateExpirationTime(AccessToken) rest_errors.RestErr
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (s *service) GetById(accessTokenId string) (*AccessToken, rest_errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token id")
	}
	accessToken, err := s.repository.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(at AccessToken) rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.repository.Create(at)
}

func (s *service) UpdateExpirationTime(at AccessToken) rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.repository.UpdateExpirationTime(at)
}
