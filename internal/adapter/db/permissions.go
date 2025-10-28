package db

import (
	"context"
)

func (db *DB) GetPermissions(ctx context.Context, userName, domainName string) ([]byte, error) {
	query := `
	WITH user_info AS (
        SELECT
            id,
            (SELECT f_user_ci_new.ci FROM backend_auth.f_user_ci_new($1)) AS ci
        FROM backend_auth.user
        WHERE
            lower(username) = $1
            OR
            lower(username) = $1 || '@' || $2
        LIMIT 1
    ),
    user_groups AS (
        SELECT
            id,
            name
        FROM backend_auth.auth_group
        WHERE auth_group.id IN (
            SELECT group_id
            FROM backend_auth.user_groups
            WHERE customuser_id = (
                SELECT id
                FROM user_info
            )
        )
    ),
    user_routes AS (
        SELECT JSON_BUILD_OBJECT(
            'method', method,
            'path', path
        ) AS permissions
        FROM backend_auth.routes
        WHERE id IN (
            SELECT route_id
            FROM backend_auth.auth_group_routes
            WHERE group_id IN (
                SELECT id
                FROM user_groups
            )
        )
    )
    SELECT JSON_BUILD_OBJECT(
        'routes', ARRAY((SELECT permissions FROM user_routes)),
        'groups', ARRAY((SELECT name FROM user_groups)),
        'ci', CASE WHEN ci IS NOT NULL THEN ci ELSE '{}'::text[] END
    ) AS v1
    FROM user_info
	`

	var v1 []byte

	err := db.pool.QueryRow(ctx, query, userName, domainName).Scan(&v1)
	if err != nil {
		return nil, err
	}

	return v1, nil
}
