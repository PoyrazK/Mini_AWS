package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepo(mock)
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         "Test User",
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.Email, user.PasswordHash, user.Name, user.Role, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_GetByEmail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepo(mock)
	email := "test@example.com"
	id := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "name", "role", "created_at", "updated_at"}).
			AddRow(id, email, "hashed", "Test User", "user", now, now))

	user, err := repo.GetByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, id, user.ID)
	assert.Equal(t, email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepo(mock)
	id := uuid.New()
	email := "test@example.com"
	now := time.Now()

	mock.ExpectQuery("SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "name", "role", "created_at", "updated_at"}).
			AddRow(id, email, "hashed", "Test User", "user", now, now))

	user, err := repo.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, id, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepo(mock)
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "updated@example.com",
		PasswordHash: "newhash",
		Name:         "Updated User",
		Role:         "admin",
		UpdatedAt:    time.Now(),
	}

	mock.ExpectExec("UPDATE users").
		WithArgs(user.Email, user.PasswordHash, user.Name, user.Role, user.UpdatedAt, user.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewUserRepo(mock)
	id1 := uuid.New()
	id2 := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, email, password_hash, name, role, created_at, updated_at FROM users").
		WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "name", "role", "created_at", "updated_at"}).
			AddRow(id1, "u1@ex.com", "h1", "U1", "user", now, now).
			AddRow(id2, "u2@ex.com", "h2", "U2", "admin", now, now))

	users, err := repo.List(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, id1, users[0].ID)
	assert.Equal(t, id2, users[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
