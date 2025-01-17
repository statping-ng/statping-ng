package middleware

import (
    "context"
    "net/http"
    "strings"

    jwt "github.com/form3tech-oss/jwt-go" // Assurez-vous d'avoir cet import
    "github.com/statping-ng/statping-ng/types/users"
)

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := extractTokenFromRequest(r)

        if tokenString == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Provide the secret key used to sign the JWT
            return []byte(core.App.Secret), nil // core.App.Secret
        })

                if err != nil || !token.Valid {
                        http.Error(w, "Unauthorized", http.StatusUnauthorized)
                        return
                }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            user := &users.User{}
                        if username, ok := claims["username"].(string); ok {
                                user.Username = username
                        }
                        if email, ok := claims["email"].(string); ok {
                                user.Email = email
                        }
                        if adminGroups, ok := claims["groups"].([]interface{}); ok {
                                for _, group := range adminGroups {
                                        if groupString, ok := group.(string); ok {
                                                user.AdminGroups = append(user.AdminGroups, groupString)
                                        }
                                }
                        }

            ctx := context.WithValue(r.Context(), userContextKey, user)
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
    })
}


func AdminAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, ok := r.Context().Value(userContextKey).(*users.User)

        if !ok || !contains(user.AdminGroups, "admin") {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func contains(s []string, str string) bool {
    for _, v := range s {
        if v == str {
            return true
        }
    }
    return false
}

func extractTokenFromRequest(r *http.Request) string {
    authHeader := r.Header.Get("Authorization")
    if authHeader != "" {
        parts := strings.Split(authHeader, " ")
        if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
            return parts[1]
        }
    }
    return ""
}