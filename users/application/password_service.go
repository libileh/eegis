package application

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}
func (ps *PasswordService) SetPassword(plainText string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	return hash, nil
}

func (ps *PasswordService) ComparePassword(plainText string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(plainText))
}
