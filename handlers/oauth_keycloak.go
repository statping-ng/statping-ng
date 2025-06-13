// oauth_keycloak.go - Keycloak OAuth handler
//
// This handler relies on the proper configuration of the Keycloak client.
// Ensure that the Keycloak client includes:
// - A GroupToRoleMapper (Token mapper for group membership, priority 0)
// - A User Realm Role mapper (roles-mapper, priority 40)
// Also, create the `statping-admin` role and map it to the appropriate groups.
// These configurations allow the userinfo token to include a roles array with 'statping-admin'.

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/statping-ng/statping-ng/types/core"
	"github.com/statping-ng/statping-ng/types/errors"
	"golang.org/x/oauth2"
)

type keycloakUserInfo struct {
	Username string   `json:"preferred_username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

func keycloakOAuth(r *http.Request) (*oAuth, error) {
	auth := core.App.OAuth
	code := r.URL.Query().Get("code")
	if code == "" {
		return nil, errors.New("code not found")
	}
	scopes := strings.Split(auth.KeycloakScopes, ",")

	conf := &oauth2.Config{
		ClientID:     core.App.OAuth.KeycloakClientID,
		ClientSecret: core.App.OAuth.KeycloakClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  auth.KeycloakEndpointAuth,
			TokenURL: auth.KeycloakEndpointToken,
		},
		RedirectURL: core.App.Domain + basePath + "oauth/keycloak",
		Scopes:      scopes,
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange token")
	}

	client := conf.Client(context.Background(), token)
	resp, err := client.Get(core.App.OAuth.KeycloakEndpointUserinfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user info")
	}
	defer resp.Body.Close()

	var userInfo keycloakUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, errors.Wrap(err, "failed to decode user info")
	}

	isAdmin := false

	// Check if the user has the 'admin' role
	for _, role := range userInfo.Roles {
		if role == "statping-admin" {
			isAdmin = true
			break
		}
	}
	log.Infoln("Keycloak user admin role:", isAdmin)

	return &oAuth{
		Email:    userInfo.Email,
		Username: userInfo.Username,
		Token:    token,
		Admin:    isAdmin,
	}, nil
}