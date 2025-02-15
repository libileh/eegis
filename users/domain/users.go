package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/libileh/eegis/users/interfaces"
	"time"

	"github.com/google/uuid"
)

type UserInvitation struct {
	Token  string    `json:"token"`
	UserId uuid.UUID `json:"user_id"`
	Expiry time.Time `json:"expiry"`
}

// swagger:model User
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	IsActive     bool      `json:"is_active" validate:"required"`
	RoleID       int16     `json:"role_id"`
}

type Password struct {
	Text *string
	Hash []byte
	User *User
}

func (u *User) NewPassword() *Password {
	return &Password{
		User: u,
	}
}

func (p *Password) Set(text string, passwordService interfaces.PasswordService) error {
	hash, err := passwordService.SetPassword(text)
	if err != nil {
		return err
	}

	p.Text = &text
	p.Hash = hash
	if p.User != nil {
		p.User.PasswordHash = hash
	}
	return nil
}

/*
GenerateToken creates a secure token for a user invitation.
It generates a UUID, hashes it using SHA-256, and returns the hash as a hexadecimal string.
*/
func (ui UserInvitation) GenerateToken() string {
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	return hex.EncodeToString(hash[:])
}
