package corbel

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type iamOauthTokenResponse struct {
	AccessToken  string `json:"accessToken,omitempty"`
	ExpiresAt    int64  `json:"expiresAt,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// OauthToken gets an access token
//
// API Docs: http://docs.silkroadiam.apiary.io/#reference/authorization/oauthtoken
func (i *IAMService) OauthToken() error {
	return i.OauthTokenBasicAuth("", "")
}

// OauthTokenBasicAuth gets an access token using username/password scheme (basic auth)
//
// API Docs: http://docs.silkroadiam.apiary.io/#reference/authorization/oauthtoken
func (i *IAMService) OauthTokenBasicAuth(username, password string) error {
	var (
		iamResponse iamOauthTokenResponse
		duration    time.Duration
	)
	signingMethod := jwt.GetSigningMethod(i.client.ClientJWTSigningMethod)
	token := jwt.New(signingMethod)
	// Required JWT Claims for SR
	token.Claims["aud"] = "http://iam.bqws.io"
	// convert to time.Duration
	duration = time.Duration(i.client.TokenExpirationTime) * time.Millisecond
	token.Claims["exp"] = time.Now().Add(duration).Unix()
	token.Claims["iss"] = i.client.ClientID
	token.Claims["scope"] = i.client.ClientScopes
	token.Claims["domain"] = i.client.ClientDomain
	token.Claims["name"] = i.client.ClientName
	// looking for basic auth pair
	if username != "" {
		token.Claims["basic_auth.username"] = username
	}
	if password != "" {
		token.Claims["basic_auth.password"] = password
	}
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(i.client.ClientSecret))
	if err != nil {
		return errJWTEncodingError
	}

	values := url.Values{}
	values.Set("grant_type", grantType)
	values.Set("assertion", tokenString)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s", i.client.URLFor("iam", "/v1.0/oauth/token")), bytes.NewBufferString(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", i.client.UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = returnErrorHTTPInterface(i.client, req, err, &iamResponse, 200)
	if err != nil {
		return err
	}

	i.client.CurrentToken = iamResponse.AccessToken
	i.client.CurrentTokenExpiresAt = iamResponse.ExpiresAt
	i.client.CurrentRefreshToken = iamResponse.RefreshToken

	return nil
}

// OauthTokenUpgrade upgrade the token using the token generated by the module Assets
//   on /assets/access and adds the scopes assigned at assets level to the current
//   logged user returning a new token with those additional scopes.
//
// API Docs: http://docs.silkroadiam.apiary.io/#reference/authorization/oauthtokenupgrade
func (i *IAMService) OauthTokenUpgrade(assetsToken string) error {
	var (
		err    error
		req    *http.Request
		res    *http.Response
		values = url.Values{}
	)
	//values := url.Values{}
	values.Set("grant_type", grantType)
	values.Set("assertion", assetsToken)
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s", i.client.URLFor("iam", "/v1.0/oauth/token/upgrade")),
		bytes.NewBufferString(values.Encode()))

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err = i.client.httpClient.Do(req)
	if res.StatusCode == 401 {
		return errHTTPNotAuthorized
	}
	return err
}
