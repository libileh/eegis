package mapper

import (
	"github.com/libileh/eegis/users/domain"
	"github.com/libileh/eegis/users/internal/infra/persistence/database"
)

func DbToUser(dbUser database.User) *domain.User {
	return &domain.User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		IsActive:  dbUser.IsActive,
		RoleID:    dbUser.RoleID,
	}
}
