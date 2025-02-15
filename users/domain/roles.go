package domain

import "github.com/libileh/eegis/users/internal/infra/persistence/database"

type Role struct {
	ID          int16  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int16  `json:"level"`
}

func DbToRole(dbRole database.UsersRole) *Role {
	return &Role{
		ID:          dbRole.ID,
		Name:        dbRole.Name,
		Description: dbRole.Description,
		Level:       dbRole.Level,
	}
}
