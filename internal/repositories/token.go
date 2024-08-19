package repositories

import (
	"database/sql"
	"time"

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

func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return RefreshTokenRepository{db: db}

}

func NewRefreshTokenRepositoryFromDataSource(dataSourceName string) RefreshTokenRepository {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err := PingWithTimeout(db, time.Second*5); err != nil {
		panic(err)
	}
	return RefreshTokenRepository{db: db}
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

func (r *RefreshTokenRepository) Shutdown() {
	r.db.Close()
}
