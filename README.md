# fasterdog

## Установка

```bash
go get github.com/SijaBakh/fasterdog@v0.1.0
```

## Использование

### 1. TokenMiddleware

Создание и использование TokenMiddleware может варьироваться в зависимости от сборки проекта.
Один из вариантов представлен ниже.

В конфигурации приложения должны быть переменные:
```yaml
middleware:
  redis_db_dsn: "redis://"
  redis_max_pool: 5
  token_secret_key: "1234567890abc"
  token_encode_algorithm: "32"
  db_auth_dsn: "postgres://"
  domain_name: "dev"
```
- `redis_db_dsn` - Подключение к редис для использования ролевой модели
- `redis_max_pool` - Количество подключений к редис для использования ролевой модели
- `token_secret_key` - Ключ шифрования токена
- `token_encode_algorithm` - Алгоритм шифрования токена
- `db_auth_dsn` - URL для подключения к БД, где хранится информация о ролевой модели
- `domain_name` - Домен необходим при выгрузке прав из БД

```go
import (
    "context"

    fasterdog "github.com/SijaBakh/fasterdog/pkg/middlewares"

    "github.com/go-chi/chi/v5"
)
// Создаём контекст приложения
ctx := context.Background()

// Инициализация мидлвара. Передаём необходимые параметры из конфигурации
mw := fasterdog.New(
		ctx,
		cfg.Middleware.RedisDBDsn,
		cfg.Middleware.DBAuthDSN,
		cfg.Middleware.TokenSecretKey,
		cfg.Middleware.TokenEncodeAlgorithm,
		cfg.Middleware.DomainName,
		cfg.Middleware.RedisMaxPool,
	)

// Регистрируем мидлвар в маршрутизаторе chi
r := chi.NewRouter()

r.Use(mw)
```

### 2. Routes

Один из примеров использования функционала записи данных эндпоинтов в БД для ролевой модели

```go
import (
    "context"

    fasterdog "github.com/SijaBakh/fasterdog/pkg/routes"

    "github.com/go-chi/chi/v5"
)

router := chi.NewRouter()
h := handlers.New()
router.Route("/api/v1/", func(r chi.Router) {
	r.Get("/service", h.Get)
	r.Post("/service", h.Post)
	r.Patch("/service", h.Patch)
	r.Delete("/service", h.Delete)
})

// Получаем все необходимые данные эндпоинтов маршрутизатора chi
routes, err := fasterdog.GetRoutes(router)
if err != nil {
	panic(err)
}

// Проверяем есть ли данные маршрутизатора в БД, если их нет или не полностью, 
// то записываем их полностью или недостающую часть
err = fasterdog.CheckRoutes(envConfig.Middleware.DBAuthDSN, routes)
if err != nil {
	panic(err)
}
```
