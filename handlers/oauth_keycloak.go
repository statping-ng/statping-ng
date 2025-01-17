package handlers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/statping-ng/statping-ng/types/core"
    "github.com/statping-ng/statping-ng/types/errors"
    "golang.org/x/oauth2"
)

func keycloakOAuth(r *http.Request) (*oAuth, error) {
    auth := core.App.OAuth
    code := r.URL.Query().Get("code")
    scopes := strings.Split(auth.KeycloakScopes, ",")
	admin_groups := strings.Split(auth.KeycloakAdminGroups, ",")

    config := &oauth2.Config{
        ClientID:     auth.KeycloakClientID,
        ClientSecret: auth.KeycloakClientSecret,
        Endpoint: oauth2.Endpoint{
            AuthURL:  auth.KeycloakEndpointAuth,
            TokenURL: auth.KeycloakEndpointToken,
        },
        RedirectURL: core.App.Domain + basePath + "oauth/keycloak",
        Scopes:      scopes,
		Admin_Groups:      admin_groups,
    }

    token, err := config.Exchange(context.Background(), code)
    if err != nil {
        return nil, fmt.Errorf("code exchange failed: %w", err)
    }

    if !token.Valid() {
        return nil, errors.New("invalid token received from Keycloak")
    }

    userInfo, err := getKeycloakUserInfo(token.AccessToken,auth)
    if err != nil {
        return nil, fmt.Errorf("failed to get user info: %w", err)
    }

    return &oAuth{
        Email:    userInfo.Email,
        Username: userInfo.Name,
        Token:    token,
    }, nil
}

type KeycloakUserInfo struct {
    Sub               string   `json:"sub"`
    EmailVerified     bool     `json:"email_verified"`
    Name              string   `json:"name"`
    PreferredUsername string   `json:"preferred_username"`
    GivenName         string   `json:"given_name"`
    FamilyName        string   `json:"family_name"`
    Email             string   `json:"email"`
    Admin_Groups            []string `json:"admin_groups"`
}

func getKeycloakUserInfo(accessToken string, auth core.OAuth) (*KeycloakUserInfo, error) {
    client := &http.Client{}

        req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, auth.KeycloakEndpointUserInfo, nil)
        if err != nil {
                return nil, fmt.Errorf("failed to create user info request: %w", err)
        }

        req.Header.Add("Authorization", "Bearer "+accessToken)

        resp, err := client.Do(req)
        if err != nil {
                return nil, fmt.Errorf("failed to get user info: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return nil, fmt.Errorf("user info request returned non-200 status: %d", resp.StatusCode)
        }

        var userInfo KeycloakUserInfo
        if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
                return nil, fmt.Errorf("failed to decode user info JSON: %w", err)
        }

        return &userInfo, nil
}