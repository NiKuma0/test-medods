package repositories

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserRepository struct {
	db *sql.DB
}

type User struct {
	UserId string
	Email  string
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func NewUserRepositoryFromDataSource(dataSourceName string) *UserRepository {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err := PingWithTimeout(db, time.Second*5); err != nil {
		panic(err)
	}
	return &UserRepository{db: db}
}

func (r *UserRepository) Get(userId string) (user User, err error) {
	return user, r.db.QueryRow(
		`
		SELECT id, email FROM users
		WHERE id=$1::UUID
		LIMIT 1
		`,
		userId,
	).Scan(&user.UserId, &user.Email)
}

func (r *UserRepository) IsExists(userId string) (isExists bool, err error) {
	return isExists, r.db.QueryRow(
		`
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE id=$1::UUID
			LIMIT 1
		)
		`,
		userId,
	).Scan(&isExists)
}

func (r *UserRepository) Shutdown() {
	r.db.Close()
}
