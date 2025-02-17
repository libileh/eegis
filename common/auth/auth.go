package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type AuthRequest struct {
	// The email of the user, required and must be a valid email format
	Email string `json:"email" validate:"required,email"`
	// The password of the user, required and between 8 to 72 characters
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type CtxRole struct {
	Value       int16
	Description string
}

var (
	ADMIN     = CtxRole{Value: 3, Description: "admin"}
	Moderator = CtxRole{Value: 2, Description: "moderator"}
	User      = CtxRole{Value: 1, Description: "user"}
)

type CtxUser struct {
	ID          uuid.UUID
	ContextRole CtxRole
}

func MapToCtxRole(roleValue float64) (CtxRole, error) {
	// Map the role value to a CtxRole
	var ctxRole CtxRole

	switch int16(roleValue) {
	case ADMIN.Value:
		ctxRole = ADMIN
	case Moderator.Value:
		ctxRole = Moderator
	case User.Value:
		ctxRole = User
	default:
		return CtxRole{}, fmt.Errorf("invalid role value: %v", roleValue)
	}
	return ctxRole, nil
}
