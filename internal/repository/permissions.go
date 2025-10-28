package repository

import (
	"context"
	"encoding/json"

	"github.com/SijaBakh/fasterdog/internal/models"
)

func (fr *FasterdogRepository) GetPermissions(ctx context.Context, userName, domainName string) (
	*models.PermissionsResult,
	error,
) {
	v1, err := fr.db.GetPermissions(ctx, userName, domainName)
	if err != nil {
		return nil, err
	}

	var permissions models.PermissionsResult
	if err := json.Unmarshal(v1, &permissions); err != nil {
		return nil, err
	}

	return &permissions, nil
}
