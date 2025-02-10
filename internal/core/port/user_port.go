package port

type UserPort interface {
	ValidateToken(token string, ownerId string) (bool, error)
}
