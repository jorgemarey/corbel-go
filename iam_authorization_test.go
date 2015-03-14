package corbel

import (
	"strings"
	"testing"
)

func TestIAMOauthToken(t *testing.T) {
	var (
		client *Client
		err    error
	)

	client, err = NewClientForEnvironment(
		nil,
		"qa",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		10)

	err = client.IAM.OauthToken()
	if got := err; got != nil {
		t.Errorf("GetToken must not fail. Got: %v  Want: nil", got)
	}

	if got, want := strings.Count(client.CurrentToken, "."), 2; got != want {
		t.Errorf("client.CurrentToken must return a token with 2 dots. Got: %v  Want: %v", got, want)
	}
}

func TestIAMOauthTokenUpgrade(t *testing.T) {
	var (
		client *Client
		err    error
	)

	client, err = NewClientForEnvironment(
		nil,
		"qa",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		10)

	err = client.IAM.OauthTokenUpgrade("aaaaaa")
	if err != errHTTPNotAuthorized {
		t.Errorf("OauthTokenUpgrade must fail since it got an invalid token. %s", err)
	}

	// // TODO: correct this test with a valid token from assets
	// if got := client.IAM.OauthTokenUpgrade("change to the assets token before it works"); got != nil {
	// 	t.Errorf("OauthTokenUpgrade failed. Got: %v Want: nil", got)
	// }

}

func TestIAMOauthTokenBasicAuth(t *testing.T) {
	var (
		client *Client
		err    error
	)

	client, err = NewClientForEnvironment(
		nil,
		"qa",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		10)

	err = client.IAM.OauthTokenBasicAuth("username", "password")
	if err != nil {
		t.Errorf("OauthTokenBasicAuth must not fail if client is correct. %s", err)
	}

	if got, want := client.CurrentToken, ""; got != want {
		t.Errorf("OauthTokenBasicAuth must not fill CurrentToken if user/password does not exists.")
	}

	// This test is done in iam_user_test.go withing the user creation workflow

}
