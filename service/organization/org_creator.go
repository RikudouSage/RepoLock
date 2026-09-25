package organization

import (
	"context"
	"fmt"

	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/entity"
	"go.chrastecky.dev/repolock/manager"
)

type Creator interface {
	CreateWithUser(ctx context.Context, organization *entity.Organization, user *entity.User) error
}

func NewCreator(
	orgMembershipManager manager.OrganizationMembership,
	orgManager manager.Organization,
	transactionCreator db.TransactionCreator,
) Creator {
	return &creator{
		orgMembershipManager: orgMembershipManager,
		orgManager:           orgManager,
		transactionCreator:   transactionCreator,
	}
}

type creator struct {
	orgMembershipManager manager.OrganizationMembership
	orgManager           manager.Organization
	transactionCreator   db.TransactionCreator
}

func (receiver *creator) CreateWithUser(ctx context.Context, organization *entity.Organization, user *entity.User) error {
	tx, err := receiver.transactionCreator.Create()
	if err != nil {
		return fmt.Errorf("failed creating transaction: %w", err)
	}
	defer tx.Rollback()

	ctx = db.WithTransaction(ctx, tx)

	if err = receiver.orgManager.Create(ctx, organization); err != nil {
		return fmt.Errorf("failed creating organization: %w", err)
	}

	if err = receiver.orgMembershipManager.Create(ctx, &entity.OrganizationMembership{
		OrganizationID: organization.ID,
		UserID:         user.ID,
		Permission:     entity.PermissionAdmin,
		Approved:       true,
	}); err != nil {
		return fmt.Errorf("failed creating organization membership: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed committing transaction: %w", err)
	}

	return nil
}
