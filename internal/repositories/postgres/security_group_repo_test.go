package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	theclouderrors "github.com/poyrazk/thecloud/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestSecurityGroupRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		sg := &domain.SecurityGroup{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			VPCID:       uuid.New(),
			Name:        "test-sg",
			Description: "desc",
			ARN:         "arn",
			CreatedAt:   time.Now(),
		}

		mock.ExpectExec("INSERT INTO security_groups").
			WithArgs(sg.ID, sg.UserID, sg.VPCID, sg.Name, sg.Description, sg.ARN, sg.CreatedAt).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = repo.Create(context.Background(), sg)
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		sg := &domain.SecurityGroup{
			ID: uuid.New(),
		}

		mock.ExpectExec("INSERT INTO security_groups").
			WillReturnError(errors.New("db error"))

		err = repo.Create(context.Background(), sg)
		assert.Error(t, err)
	})
}

func TestSecurityGroupRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)
		now := time.Now()

		mock.ExpectQuery("SELECT id, user_id, vpc_id, name, description, arn, created_at FROM security_groups").
			WithArgs(id, userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "vpc_id", "name", "description", "arn", "created_at"}).
				AddRow(id, userID, uuid.New(), "test-sg", "desc", "arn", now))

		mock.ExpectQuery("SELECT id, group_id, direction, protocol, port_min, port_max, cidr, priority, created_at FROM security_rules").
			WithArgs(id).
			WillReturnRows(pgxmock.NewRows([]string{"id", "group_id", "direction", "protocol", "port_min", "port_max", "cidr", "priority", "created_at"}).
				AddRow(uuid.New(), id, string(domain.RuleIngress), "tcp", 80, 80, "0.0.0.0/0", 100, now))

		sg, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		if sg != nil {
			assert.Len(t, sg.Rules, 1)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectQuery("SELECT id, user_id, vpc_id, name, description, arn, created_at FROM security_groups").
			WithArgs(id, userID).
			WillReturnError(pgx.ErrNoRows)

		sg, err := repo.GetByID(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, sg)
		theCloudErr, ok := err.(*theclouderrors.Error)
		if ok {
			assert.Equal(t, theclouderrors.NotFound, theCloudErr.Type)
		}
	})
}

func TestSecurityGroupRepository_GetByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		vpcID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)
		now := time.Now()
		name := "test-sg"

		mock.ExpectQuery("SELECT id, user_id, vpc_id, name, description, arn, created_at FROM security_groups").
			WithArgs(vpcID, name, userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "vpc_id", "name", "description", "arn", "created_at"}).
				AddRow(id, userID, vpcID, name, "desc", "arn", now))

		mock.ExpectQuery("SELECT id, group_id, direction, protocol, port_min, port_max, cidr, priority, created_at FROM security_rules").
			WithArgs(id).
			WillReturnRows(pgxmock.NewRows([]string{"id", "group_id", "direction", "protocol", "port_min", "port_max", "cidr", "priority", "created_at"}).
				AddRow(uuid.New(), id, string(domain.RuleIngress), "tcp", 80, 80, "0.0.0.0/0", 100, now))

		sg, err := repo.GetByName(ctx, vpcID, name)
		assert.NoError(t, err)
		assert.NotNil(t, sg)
		assert.Equal(t, id, sg.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		userID := uuid.New()
		vpcID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)
		name := "test-sg"

		mock.ExpectQuery("SELECT id, user_id, vpc_id, name, description, arn, created_at FROM security_groups").
			WithArgs(vpcID, name, userID).
			WillReturnError(pgx.ErrNoRows)

		sg, err := repo.GetByName(ctx, vpcID, name)
		assert.Error(t, err)
		assert.Nil(t, sg)
		theCloudErr, ok := err.(*theclouderrors.Error)
		if ok {
			assert.Equal(t, theclouderrors.NotFound, theCloudErr.Type)
		}
	})
}

func TestSecurityGroupRepository_ListByVPC(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		vpcID := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)
		now := time.Now()

		mock.ExpectQuery("SELECT id, user_id, vpc_id, name, description, arn, created_at FROM security_groups").
			WithArgs(vpcID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "vpc_id", "name", "description", "arn", "created_at"}).
				AddRow(uuid.New(), userID, vpcID, "test-sg", "desc", "arn", now))

		groups, err := repo.ListByVPC(ctx, vpcID)
		assert.NoError(t, err)
		assert.Len(t, groups, 1)
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		vpcID := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectQuery("SELECT id, user_id, vpc_id, name, description, arn, created_at FROM security_groups").
			WithArgs(vpcID, userID).
			WillReturnError(errors.New("db error"))

		groups, err := repo.ListByVPC(ctx, vpcID)
		assert.Error(t, err)
		assert.Nil(t, groups)
	})
}

func TestSecurityGroupRepository_AddRule(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		rule := &domain.SecurityRule{
			ID:        uuid.New(),
			GroupID:   uuid.New(),
			Direction: domain.RuleIngress,
			Protocol:  "tcp",
			PortMin:   80,
			PortMax:   80,
			CIDR:      "0.0.0.0/0",
			Priority:  100,
			CreatedAt: time.Now(),
		}

		mock.ExpectExec("INSERT INTO security_rules").
			WithArgs(rule.ID, rule.GroupID, rule.Direction, rule.Protocol, rule.PortMin, rule.PortMax, rule.CIDR, rule.Priority, rule.CreatedAt).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = repo.AddRule(context.Background(), rule)
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		rule := &domain.SecurityRule{
			ID: uuid.New(),
		}

		mock.ExpectExec("INSERT INTO security_rules").
			WillReturnError(errors.New("db error"))

		err = repo.AddRule(context.Background(), rule)
		assert.Error(t, err)
	})
}

func TestSecurityGroupRepository_DeleteRule(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		ruleID := uuid.New()

		mock.ExpectExec("DELETE FROM security_rules").
			WithArgs(ruleID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err = repo.DeleteRule(context.Background(), ruleID)
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		ruleID := uuid.New()

		mock.ExpectExec("DELETE FROM security_rules").
			WithArgs(ruleID).
			WillReturnError(errors.New("db error"))

		err = repo.DeleteRule(context.Background(), ruleID)
		assert.Error(t, err)
	})
}

func TestSecurityGroupRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectExec("DELETE FROM security_groups").
			WithArgs(id, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err = repo.Delete(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectExec("DELETE FROM security_groups").
			WithArgs(id, userID).
			WillReturnError(errors.New("db error"))

		err = repo.Delete(ctx, id)
		assert.Error(t, err)
	})
}

func TestSecurityGroupRepository_AddInstanceToGroup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		instanceID := uuid.New()
		groupID := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT user_id FROM instances WHERE id = \\$1").
			WithArgs(instanceID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(userID))
		mock.ExpectQuery("SELECT user_id FROM security_groups WHERE id = \\$1").
			WithArgs(groupID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(userID))
		mock.ExpectExec("INSERT INTO instance_security_groups").
			WithArgs(instanceID, groupID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mock.ExpectCommit()

		err = repo.AddInstanceToGroup(ctx, instanceID, groupID)
		assert.NoError(t, err)
	})

	t.Run("instance not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		instanceID := uuid.New()
		groupID := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT user_id FROM instances WHERE id = \\$1").
			WithArgs(instanceID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectRollback()

		err = repo.AddInstanceToGroup(ctx, instanceID, groupID)
		assert.Error(t, err)
		theCloudErr, ok := err.(*theclouderrors.Error)
		if ok {
			assert.Equal(t, theclouderrors.NotFound, theCloudErr.Type)
		}
	})

	t.Run("security group not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		instanceID := uuid.New()
		groupID := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT user_id FROM instances WHERE id = \\$1").
			WithArgs(instanceID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(userID))
		mock.ExpectQuery("SELECT user_id FROM security_groups WHERE id = \\$1").
			WithArgs(groupID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectRollback()

		err = repo.AddInstanceToGroup(ctx, instanceID, groupID)
		assert.Error(t, err)
		theCloudErr, ok := err.(*theclouderrors.Error)
		if ok {
			assert.Equal(t, theclouderrors.NotFound, theCloudErr.Type)
		}
	})
}

func TestSecurityGroupRepository_RemoveInstanceFromGroup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		instanceID := uuid.New()
		groupID := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT user_id FROM instances WHERE id = \\$1").
			WithArgs(instanceID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(userID))
		mock.ExpectQuery("SELECT user_id FROM security_groups WHERE id = \\$1").
			WithArgs(groupID).
			WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(userID))
		mock.ExpectExec("DELETE FROM instance_security_groups").
			WithArgs(instanceID, groupID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()

		err = repo.RemoveInstanceFromGroup(ctx, instanceID, groupID)
		assert.NoError(t, err)
	})

	t.Run("instance not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		instanceID := uuid.New()
		groupID := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT user_id FROM instances WHERE id = \\$1").
			WithArgs(instanceID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectRollback()

		err = repo.RemoveInstanceFromGroup(ctx, instanceID, groupID)
		assert.Error(t, err)
		theCloudErr, ok := err.(*theclouderrors.Error)
		if ok {
			assert.Equal(t, theclouderrors.NotFound, theCloudErr.Type)
		}
	})
}

func TestSecurityGroupRepository_ListInstanceGroups(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		instanceID := uuid.New()
		now := time.Now()

		mock.ExpectQuery("SELECT sg.id, sg.user_id, sg.vpc_id, sg.name, sg.description, sg.arn, sg.created_at").
			WithArgs(instanceID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "vpc_id", "name", "description", "arn", "created_at"}).
				AddRow(uuid.New(), uuid.New(), uuid.New(), "test-sg", "desc", "arn", now))

		groups, err := repo.ListInstanceGroups(context.Background(), instanceID)
		assert.NoError(t, err)
		assert.Len(t, groups, 1)
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewSecurityGroupRepository(mock)
		instanceID := uuid.New()

		mock.ExpectQuery("SELECT sg.id, sg.user_id, sg.vpc_id, sg.name, sg.description, sg.arn, sg.created_at").
			WithArgs(instanceID).
			WillReturnError(errors.New("db error"))

		groups, err := repo.ListInstanceGroups(context.Background(), instanceID)
		assert.Error(t, err)
		assert.Nil(t, groups)
	})
}
