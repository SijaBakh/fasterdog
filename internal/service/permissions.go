package service

import (
	"context"
	"encoding/json"

	"github.com/SijaBakh/fasterdog/internal/models"
)

func (fs *FasterdogService) GetPermissions(ctx context.Context, userName, domainName string) (*models.PermissionsResult, error) {
	v1, err := fs.repo.GetPermissions(ctx, userName, domainName)
	if err != nil {
		return nil, err
	}

	var permissions models.PermissionsResult
	if err := json.Unmarshal(v1, &permissions); err != nil {
		return nil, err
	}

	return &permissions, nil
}

func (fs *FasterdogService) RGetPermissions(ctx context.Context, username string) (*models.PermissionsResult, error) {
	pr, err := fs.rClient.RGetPermissions(ctx, username)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
