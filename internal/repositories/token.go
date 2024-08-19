package repositories

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

type Token struct {
	Id     int
	Hash   []byte
	IsUsed bool
}

func NewTokenRepository(db *sql.DB) RefreshTokenRepository {
	return RefreshTokenRepository{db: db}

}

type TokenRepository interface {
	SaveRefreshToken(userId, tokenHash string) error
	DeleteRefreshToken(tokenHash string) error
	IsRefreshTokenValid(userId, tokenHash string) (bool, error)
}

func NewTokenRepositoryFromDataSource(dataSourceName string) (r UserRepository) {
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
func (r *RefreshTokenRepository) SaveRefreshToken(userId, tokenHash string) (err error) {
	_, err = r.db.Exec(
		`
		INSERT INTO tokens (user_id, hash)
		VALUES ($1::UUID, $2)
		`,
		userId, tokenHash,
	)
	return

}
func (r *RefreshTokenRepository) DeleteRefreshToken(tokenHash string) (err error) {
	_, err = r.db.Exec(
		`
		DELETE FROM tokens WHERE hash=$1
		`,
		tokenHash,
	)
	return
}

func (r *RefreshTokenRepository) IsRefreshTokenValid(userId, tokenHash string) (bool, error) {
	var count int
	row := r.db.QueryRow(
		`
		SELECT COUNT(*) FROM tokens
		WHERE user_id=$1::UUID AND hash=$2
		LIMIT 1
		`,
		userId, tokenHash,
	)
	err := row.Scan(&count)
	return count > 0, err
}

func (r *RefreshTokenRepository) Get(userId string) (token Token, isExists bool, err error) {
	rows, err := r.db.Query(
		`
		SELECT id, hash, is_used FROM tokens
		WHERE user_id=$1::UUID
		LIMIT 1
		`,
		userId,
	)
	isExists = rows.Next()
	if !isExists {
		return
	}
	err = rows.Scan(
		&token.Id,
		&token.Hash,
		&token.IsUsed,
	)
	return
}

func (r *RefreshTokenRepository) Shutdown() {
	r.db.Close()
}
