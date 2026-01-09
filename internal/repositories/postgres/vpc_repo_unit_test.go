package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestVpcRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVpcRepository(mock)
	vpc := &domain.VPC{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Name:      "test-vpc",
		CIDRBlock: "10.0.0.0/16",
		NetworkID: "net-1",
		VXLANID:   100,
		Status:    "available",
		ARN:       "arn",
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO vpcs").
		WithArgs(vpc.ID, vpc.UserID, vpc.Name, vpc.CIDRBlock, vpc.NetworkID, vpc.VXLANID, vpc.Status, vpc.ARN, vpc.CreatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), vpc)
	assert.NoError(t, err)
}

func TestVpcRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVpcRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, COALESCE\\(cidr_block::text, ''\\), network_id, vxlan_id, status, arn, created_at FROM vpcs").
		WithArgs(id, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "cidr_block", "network_id", "vxlan_id", "status", "arn", "created_at"}).
			AddRow(id, userID, "test-vpc", "10.0.0.0/16", "net-1", 100, "available", "arn", now))

	vpc, err := repo.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, vpc)
	assert.Equal(t, id, vpc.ID)
}

func TestVpcRepository_GetByName(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVpcRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()
	name := "test-vpc"

	mock.ExpectQuery("SELECT id, user_id, name, COALESCE\\(cidr_block::text, ''\\), network_id, vxlan_id, status, arn, created_at FROM vpcs").
		WithArgs(name, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "cidr_block", "network_id", "vxlan_id", "status", "arn", "created_at"}).
			AddRow(id, userID, name, "10.0.0.0/16", "net-1", 100, "available", "arn", now))

	vpc, err := repo.GetByName(ctx, name)
	assert.NoError(t, err)
	assert.NotNil(t, vpc)
	assert.Equal(t, id, vpc.ID)
}

func TestVpcRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVpcRepository(mock)
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, COALESCE\\(cidr_block::text, ''\\), network_id, vxlan_id, status, arn, created_at FROM vpcs").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "cidr_block", "network_id", "vxlan_id", "status", "arn", "created_at"}).
			AddRow(uuid.New(), userID, "test-vpc", "10.0.0.0/16", "net-1", 100, "available", "arn", now))

	vpcs, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, vpcs, 1)
}

func TestVpcRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVpcRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)

	mock.ExpectExec("DELETE FROM vpcs").
		WithArgs(id, userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(ctx, id)
	assert.NoError(t, err)
}
