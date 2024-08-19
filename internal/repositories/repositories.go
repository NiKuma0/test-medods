package repositories

type Repositories struct {
	User interface {
		Get(userId string) (User, error)
		IsExists(userId string) (bool, error)
	}
	Token interface {
		Get(userId string) ([]byte, error)
		Create(userId string, hash []byte)
		SetIsUsed(id int, isUsed bool) error
	}
}
