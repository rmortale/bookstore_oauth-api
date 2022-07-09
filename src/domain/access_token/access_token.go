package access_token

import (
	"fmt"
	"github.com/rmortale/bookstore_oauth-api/src/utils/crypto_utils"
	"github.com/rmortale/bookstore_oauth-api/src/utils/errors"
	"strings"
	"time"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

type AccessTokenRequest struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
}

func (at *AccessTokenRequest) Validate() *errors.RestErr {
	if at.GrantType != grantTypePassword || at.GrantType != grantTypeClientCredentials {
		return errors.NewBadRequestError("invalid grant_type parameter")
	}
	return nil
}

func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequestError("invalid access token")
	}
	if at.UserId <= 0 {
		return errors.NewBadRequestError("invalid userid")
	}
	if at.ClientId <= 0 {
		return errors.NewBadRequestError("invalid ClientId")
	}
	if at.Expires <= 0 {
		return errors.NewBadRequestError("invalid Expires")
	}
	return nil
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
}
