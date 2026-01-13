package postgres

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestNewDualDB(t *testing.T) {
	primary, _ := pgxmock.NewPool()
	replica, _ := pgxmock.NewPool()
	defer primary.Close()
	defer replica.Close()

	dual := NewDualDB(primary, replica)
	assert.NotNil(t, dual)
	assert.Equal(t, primary, dual.primary)
	assert.Equal(t, replica, dual.replica)

	// Test fallback
	dual2 := NewDualDB(primary, nil)
	assert.Equal(t, primary, dual2.primary)
	assert.Equal(t, primary, dual2.replica)
}

func TestDualDB_Operations(t *testing.T) {
	primary, _ := pgxmock.NewPool()
	replica, _ := pgxmock.NewPool()
	defer primary.Close()
	defer replica.Close()

	dual := NewDualDB(primary, replica)
	ctx := context.Background()

	// Exec should go to primary
	primary.ExpectExec("INSERT").WillReturnResult(pgxmock.NewResult("INSERT", 1))
	_, err := dual.Exec(ctx, "INSERT INTO test DEFAULT VALUES")
	assert.NoError(t, err)

	// Query should go to replica
	replica.ExpectQuery("SELECT").WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	rows, err := dual.Query(ctx, "SELECT id FROM test")
	assert.NoError(t, err)
	rows.Close()

	// QueryRow should go to replica
	replica.ExpectQuery("SELECT").WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	row := dual.QueryRow(ctx, "SELECT id FROM test WHERE id = 1")
	var id int
	err = row.Scan(&id)
	assert.NoError(t, err)

	// Begin should go to primary
	primary.ExpectBegin()
	_, err = dual.Begin(ctx)
	assert.NoError(t, err)

	// Ping should go to primary
	primary.ExpectPing()
	err = dual.Ping(ctx)
	assert.NoError(t, err)

	// Close should close both
	primary.ExpectClose()
	replica.ExpectClose()
	dual.Close()
}

func TestDualDB_CloseSame(t *testing.T) {
	primary, _ := pgxmock.NewPool()
	defer primary.Close()

	dual := NewDualDB(primary, nil)
	primary.ExpectClose()
	dual.Close()
}

