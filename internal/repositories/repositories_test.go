package repositories_test

import (
	"database/sql"
	"testing"

	"jwt-service/internal/repositories"

	"github.com/DATA-DOG/go-txdb"
	"github.com/stretchr/testify/assert"
)

func init() {
	txdb.Register("pgx_txdb", "pgx", "postgres://postgres:postgres@localhost:5432/app")
}

func TestRefreshTokenRepository(t *testing.T) {
	db, err := sql.Open("pgx_txdb", "")
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewRefreshTokenRepository(db)

	t.Run("SaveRefreshToken", func(t *testing.T) {
		userId := "550e8400-e29b-41d4-a716-446655440000"
		tokenHash := "abc123"
		err := repo.SaveRefreshToken(userId, tokenHash)
		assert.NoError(t, err)

		isValid, err := repo.IsRefreshTokenValid(userId, tokenHash)
		assert.NoError(t, err)
		assert.True(t, isValid)
	})

	t.Run("DeleteRefreshToken", func(t *testing.T) {
		userId := "550e8400-e29b-41d4-a716-446655440001"
		tokenHash := "hash2"
		err := repo.DeleteRefreshToken(tokenHash)
		assert.NoError(t, err)

		isValid, err := repo.IsRefreshTokenValid(userId, tokenHash)
		assert.NoError(t, err)
		assert.False(t, isValid)
	})

	t.Run("IsRefreshTokenValid", func(t *testing.T) {
		validUserId := "550e8400-e29b-41d4-a716-446655440000"
		validTokenHash := "hash1"
		invalidUserId := "550e8400-e29b-41d4-a716-446655440002"
		invalidTokenHash := "hash4"

		isValid, err := repo.IsRefreshTokenValid(validUserId, validTokenHash)
		assert.NoError(t, err)
		assert.True(t, isValid)

		isValid, err = repo.IsRefreshTokenValid(invalidUserId, invalidTokenHash)
		assert.NoError(t, err)
		assert.False(t, isValid)
	})
}

func TestUserRepository(t *testing.T) {
	db, err := sql.Open("pgx_txdb", "")
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	t.Run("Get", func(t *testing.T) {
		userId := "550e8400-e29b-41d4-a716-446655440000"
		user, err := repo.Get(userId)
		assert.NoError(t, err)
		assert.Equal(t, userId, user.UserId)
		assert.Equal(t, "test1@example.com", user.Email)
	})

	t.Run("IsExists", func(t *testing.T) {
		existingUserId := "550e8400-e29b-41d4-a716-446655440000"
		nonExistingUserId := "550e8400-e29b-41d4-a716-446655440003"

		isExists, err := repo.IsExists(existingUserId)
		assert.NoError(t, err)
		assert.True(t, isExists)

		isExists, err = repo.IsExists(nonExistingUserId)
		assert.NoError(t, err)
		assert.False(t, isExists)
	})
}
