package repositories

import (
	"database/sql"

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

func NewUserRepositoryFromDataSource(dataSourceName string) (r UserRepository) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	r = UserRepository{db: db}
	return
}

func (r *UserRepository) Get(userId string) (user User, err error) {
	err = r.db.QueryRow(
		`
		SELECT id, email FROM users
		WHERE id=$1::UUID
		LIMIT 1
		`,
		userId,
	).Scan(&user.UserId, &user.Email)
	return
}

func (r *UserRepository) IsExists(userId string) (isExists bool, err error) {
	err = r.db.QueryRow(
		`
		SELECT EXISTS (
			SELECT * FROM users
			WHERE id=$1::UUID
			LIMIT 1
		)
		`,
		userId,
	).Scan(&isExists)
	return
}

func (r *UserRepository) Shutdown() {
	r.db.Close()
}
