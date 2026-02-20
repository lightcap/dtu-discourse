// Package middleware provides HTTP middleware for the DTU server.
// Authentication follows Discourse's current requirement: credentials
// must be passed via HTTP headers (Api-Key + Api-Username), not query
// params. This matches the post-April-2020 Discourse authentication model.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/lightcap/dtu-discourse/internal/store"
)

type contextKey string

const (
	ContextKeyUsername contextKey = "api_username"
	ContextKeyIsAdmin contextKey = "is_admin"
)

func Auth(s *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// SSO browser redirects never carry API key headers
			p := strings.TrimSuffix(r.URL.Path, "/")
			if p == "/session/sso" || p == "/session/sso_login" || p == "/session/sso_provider" {
				next.ServeHTTP(w, r)
				return
			}

			apiKey := r.Header.Get("Api-Key")
			apiUsername := r.Header.Get("Api-Username")

			// Fall back to query parameters (used by some SDK clients)
			if apiKey == "" {
				apiKey = r.URL.Query().Get("api_key")
			}
			if apiUsername == "" {
				apiUsername = r.URL.Query().Get("api_username")
			}

			if apiKey == "" {
				http.Error(w, `{"errors":["not logged in"],"error_type":"not_logged_in"}`, http.StatusForbidden)
				return
			}

			keyOwner, valid := s.ValidateAPIKey(apiKey)
			if !valid {
				http.Error(w, `{"errors":["invalid api key"],"error_type":"invalid_access"}`, http.StatusForbidden)
				return
			}

			// If Api-Username is not set, use the key owner
			if apiUsername == "" {
				apiUsername = keyOwner
			}

			u := s.GetUserByUsername(apiUsername)
			isAdmin := u != nil && u.Admin

			ctx := context.WithValue(r.Context(), ContextKeyUsername, apiUsername)
			ctx = context.WithValue(ctx, ContextKeyIsAdmin, isAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUsername(r *http.Request) string {
	if v, ok := r.Context().Value(ContextKeyUsername).(string); ok {
		return v
	}
	return ""
}

func IsAdmin(r *http.Request) bool {
	if v, ok := r.Context().Value(ContextKeyIsAdmin).(bool); ok {
		return v
	}
	return false
}
