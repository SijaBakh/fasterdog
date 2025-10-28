package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"fasterdog/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type ctxKey string

const (
	CtxUsernameKey ctxKey = "username"
	CtxGroupsKey   ctxKey = "groups"
	CtxCIKey       ctxKey = "ci"
	CtxAdminISKey  ctxKey = "admin_is"
)

func (m *TokenMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if m.isExcludedPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("x-interservice-token")
		if token == "" {
			http.Error(w, "Missing token!", http.StatusUnauthorized)
			return
		}

		tokenData, err := m.decodeToken(token, []byte(m.TokenSecretKey), m.TokenEncodeAlgorithm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		b, _ := json.Marshal(tokenData)
		var c models.Claims
		_ = json.Unmarshal(b, &c)

		if len(c.Groups) == 0 {
			slog.Info(
				"У пользователя нет групп доступа (роли в WD)",
				"username", c.Username,
			)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUsernameKey, c.Username)
		ctx = context.WithValue(ctx, CtxGroupsKey, c.Groups)

		if slices.Contains(c.Groups, "Администраторы WD") {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		permissions, err := m.RedisClient.GetPermissions(r.Context(), c.Username)
		if permissions == nil && errors.Is(err, redis.Nil) {
			slog.Info(
				"У пользователя нет прав доступа",
				"username", c.Username,
			)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if err != nil {
			slog.Error(
				"При попытке получить права доступа из redis возникла ошибка",
				"error", err.Error(),
				"username", c.Username,
			)
			permissions, err = m.FasterdogRepository.GetPermissions(r.Context(), c.Username, m.DomainName)
			if permissions == nil && err == nil {
				slog.Info(
					"У пользователя нет прав доступа",
					"username", c.Username,
				)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if err != nil {
				slog.Error(
					"При попытке получить права доступа из PG возникла ошибка",
					"error", err.Error(),
					"username", c.Username,
				)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		route := models.RoutePermission{
			Method: r.Method,
			Path:   strings.Split(r.URL.Path, "?")[0],
		}

		if !isRouteAllowed(route, permissions.Routes) {
			slog.Info(
				"У пользователя нет прав доступа к эндпоинту",
				"username", c.Username,
				"route", route,
			)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		for _, g := range c.Groups {
			if slices.Contains(m.CiGroups, g) {
				ctx = context.WithValue(ctx, CtxCIKey, permissions.CI)
				ctx = context.WithValue(ctx, CtxAdminISKey, !slices.Contains(c.Groups, "Администраторы УАБД"))
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m TokenMiddleware) isExcludedPath(inputPath string) bool {
	excludedPaths := []string{
		"/ping",
		"/docs",
		"/openapi.json",
		"/redoc",
		"/static",
		"/metrics",
		"//",
	}
	for _, prefix := range excludedPaths {
		if strings.HasPrefix(inputPath, prefix) {
			return true
		}
	}
	return false
}

func (m TokenMiddleware) decodeToken(tokenStr string, secret any, algorithm string) (map[string]any, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		switch algorithm {
		case "HS256":
			secret, ok := secret.([]byte)
			if !ok {
				return nil, errors.New("unsupported secret type")
			}
			return secret, nil
		default:
			return nil, errors.New("unsupported algorithm")
		}
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if data, ok := claims["data"].(map[string]any); ok {
			return data, nil
		}
		if data, ok := claims["data"]; ok && data == nil {
			return nil, errors.New("missing data in token")
		}
		return nil, errors.New("data field missing or invalid type")
	}
	return nil, errors.New("invalid token claims")
}

func isRouteAllowed(route models.RoutePermission, allowed []models.RoutePermission) bool {
	for _, r := range allowed {
		if r.Method == route.Method && r.Path == route.Path {
			return true
		}
	}
	return false
}
