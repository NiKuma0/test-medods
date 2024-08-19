package repositories

type Repositories struct {
	User interface {
		Get(userId string) (User, error)
		IsExists(userId string) (bool, error)
	}
	Token interface {
		SaveRefreshToken(userId, tokenHash string) error
		DeleteRefreshToken(tokenHash string) error
		IsRefreshTokenValid(userId, tokenHash string) (bool, error)
	}
}
