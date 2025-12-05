package middlewares

import (
	"context"
	"net/http"

	"github.com/SijaBakh/fasterdog/internal/service"
)

type TokenMiddleware struct {
	TokenSecretKey       string
	TokenEncodeAlgorithm string
	DomainName           string
	FasterdogService     service.FasterdogServiceInterfaces
	CiGroups             []string
}

func New(
	ctx context.Context,
	redisDSN, authDSN, tokenSecretKey, tokenEncodeAlgorithm, domainName string,
	redisMP int,
) func(http.Handler) http.Handler {
	fs, err := service.New(ctx, redisDSN, authDSN, redisMP)
	if err != nil {
		panic(err)
	}

	m := &TokenMiddleware{
		TokenSecretKey:       tokenSecretKey,
		TokenEncodeAlgorithm: tokenEncodeAlgorithm,
		DomainName:           domainName,
		FasterdogService:     fs,
		CiGroups:             []string{"Администраторы ИС", "Администраторы ИС Тест"},
	}
	return m.Handler
}
