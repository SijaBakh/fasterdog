package redis

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/SijaBakh/fasterdog/internal/models"
)

func (rc *RedisClient) GetPermissions(ctx context.Context, username string) (*models.PermissionsResult, error) {
	val, err := rc.client.HGet(ctx, "wd_permissions_v1", strings.ToLower(username)).Bytes()
	if err != nil {
		return nil, err
	}

	var res models.PermissionsResult
	if err := json.Unmarshal(val, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
