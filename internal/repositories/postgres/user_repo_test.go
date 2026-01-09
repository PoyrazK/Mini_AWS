package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/poyrazk/thecloud/internal/core/domain"
	theclouderrors "github.com/poyrazk/thecloud/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)
		user := &domain.User{
			ID: uuid.New(),
		}

		mock.ExpectExec("INSERT INTO users").
			WillReturnError(errors.New("db error"))

		err = repo.Create(context.Background(), user)
		assert.Error(t, err)
	})
}

func TestUserRepo_GetByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)
		email := "test@example.com"

		mock.ExpectQuery("SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE email = \\$1").
			WithArgs(email).
			WillReturnError(pgx.ErrNoRows)

		user, err := repo.GetByEmail(context.Background(), email)
		assert.Error(t, err)
		assert.Nil(t, user)
		// Assuming repo returns custom error or native error
		// Checking if it matches domain logic. Usually repo wraps error.
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)
		email := "test@example.com"

		mock.ExpectQuery("SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE email = \\$1").
			WithArgs(email).
			WillReturnError(errors.New("db error"))

		user, err := repo.GetByEmail(context.Background(), email)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepo_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)
		id := uuid.New()

		mock.ExpectQuery("SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE id = \\$1").
			WithArgs(id).
			WillReturnError(pgx.ErrNoRows)

		user, err := repo.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, user)
		theCloudErr, ok := err.(*theclouderrors.Error)
		if ok {
			assert.Equal(t, theclouderrors.NotFound, theCloudErr.Type)
		}
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)
		id := uuid.New()

		mock.ExpectQuery("SELECT id, email, password_hash, name, role, created_at, updated_at FROM users WHERE id = \\$1").
			WithArgs(id).
			WillReturnError(errors.New("db error"))

		user, err := repo.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepo_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)
		user := &domain.User{
			ID: uuid.New(),
		}

		mock.ExpectExec("UPDATE users").
			WillReturnError(errors.New("db error"))

		err = repo.Update(context.Background(), user)
		assert.Error(t, err)
	})
}

func TestUserRepo_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)

		mock.ExpectQuery("SELECT id, email, password_hash, name, role, created_at, updated_at FROM users").
			WillReturnError(errors.New("db error"))

		users, err := repo.List(context.Background())
		assert.Error(t, err)
		assert.Nil(t, users)
	})

	t.Run("scan error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewUserRepo(mock)
		now := time.Now()

		mock.ExpectQuery("SELECT id, email, password_hash, name, role, created_at, updated_at FROM users").
			WillReturnRows(pgxmock.NewRows([]string{"id", "email", "password_hash", "name", "role", "created_at", "updated_at"}).
				AddRow("invalid-uuid", "u1@ex.com", "h1", "U1", "user", now, now))

		users, err := repo.List(context.Background())
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}
