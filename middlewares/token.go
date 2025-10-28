package middlewares

import (
	"context"
	"net/http"

	"github.com/SijaBakh/fasterdog/internal/adapter/db"
	"github.com/SijaBakh/fasterdog/internal/adapter/redis"
	"github.com/SijaBakh/fasterdog/internal/repository"
)

type TokenMiddleware struct {
	TokenSecretKey       string
	TokenEncodeAlgorithm string
	DomainName           string
	RedisClient          *redis.RedisClient
	FasterdogRepository  *repository.FasterdogRepository
	CiGroups             []string
}

func New(
	ctx context.Context,
	redisDSN, authDSN, tokenSecretKey, tokenEncodeAlgorithm, domainName string,
	redisMP int,
) func(http.Handler) http.Handler {
	rc, err := redis.New(redisDSN, redisMP)
	if err != nil {
		panic(err)
	}

	db, err := db.New(authDSN, ctx)
	if err != nil {
		panic(err)
	}
	fr := repository.New(db)
	m := &TokenMiddleware{
		TokenSecretKey:       tokenSecretKey,
		TokenEncodeAlgorithm: tokenEncodeAlgorithm,
		DomainName:           domainName,
		RedisClient:          rc,
		FasterdogRepository:  fr,
		CiGroups:             []string{"Администраторы ИС", "Администраторы ИС Тест"},
	}
	return m.Handler
}
