package interfaces

// PasswordService defines the interface for password-related operations
type PasswordService interface {
	SetPassword(plainText string) ([]byte, error)
	ComparePassword(plainText string, hash []byte) error
}
