package repository

import (
	"context"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/users/domain"
	"github.com/libileh/eegis/users/internal/infra/persistence/database"
)

type RoleStore struct {
	Queries *database.Queries
}

func (r *RoleStore) GetByRoleName(ctx context.Context, roleName string) (*domain.Role, *errors.CustomError) {
	dbRole, err := r.Queries.GetByRoleName(ctx, roleName)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	//func() string { if dbRole.Description.Valid { return dbRole.Description.String } else { return "" } }(),
	return domain.DbToRole(dbRole), nil
}

func (r *RoleStore) GetById(ctx context.Context, roleID int16) (*domain.Role, *errors.CustomError) {
	dbRole, err := r.Queries.GetById(ctx, roleID)
	if err != nil {
		return nil, errors.HandleDBError(err)
	}
	return domain.DbToRole(dbRole), nil
}
