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

func TestSubnetRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSubnetRepository(mock)
	s := &domain.Subnet{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		VPCID:            uuid.New(),
		Name:             "subnet-1",
		CIDRBlock:        "10.0.1.0/24",
		AvailabilityZone: "us-east-1a",
		GatewayIP:        "10.0.1.1",
		ARN:              "arn",
		Status:           "available",
		CreatedAt:        time.Now(),
	}

	mock.ExpectExec("INSERT INTO subnets").
		WithArgs(s.ID, s.UserID, s.VPCID, s.Name, s.CIDRBlock, s.AvailabilityZone, s.GatewayIP, s.ARN, s.Status, s.CreatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), s)
	assert.NoError(t, err)
}

func TestSubnetRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSubnetRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, vpc_id, name, cidr_block::text, availability_zone, COALESCE\\(gateway_ip::text, ''\\), arn, status, created_at FROM subnets").
		WithArgs(id, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "vpc_id", "name", "cidr_block", "availability_zone", "gateway_ip", "arn", "status", "created_at"}).
			AddRow(id, userID, uuid.New(), "subnet-1", "10.0.1.0/24", "us-east-1a", "10.0.1.1", "arn", "available", now))

	s, err := repo.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, id, s.ID)
}

func TestSubnetRepository_ListByVPC(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewSubnetRepository(mock)
	vpcID := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, vpc_id, name, cidr_block::text, availability_zone, COALESCE\\(gateway_ip::text, ''\\), arn, status, created_at FROM subnets").
		WithArgs(vpcID, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "vpc_id", "name", "cidr_block", "availability_zone", "gateway_ip", "arn", "status", "created_at"}).
			AddRow(uuid.New(), userID, vpcID, "subnet-1", "10.0.1.0/24", "us-east-1a", "10.0.1.1", "arn", "available", now))

	subnets, err := repo.ListByVPC(ctx, vpcID)
	assert.NoError(t, err)
	assert.Len(t, subnets, 1)
}
