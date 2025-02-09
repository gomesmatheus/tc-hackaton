package port

type UserPort interface {
	ValidateToken(token string) (string, error)
}
