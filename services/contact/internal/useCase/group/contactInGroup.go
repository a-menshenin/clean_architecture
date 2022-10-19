package group

import (
	"github.com/google/uuid"

	"architecture_go_2/pkg/type/context"
	"architecture_go_2/services/contact/internal/domain/contact"
)

func (uc *UseCase) CreateContactIntoGroup(ctx context.Context, groupID uuid.UUID, contacts ...*contact.Contact) ([]*contact.Contact, error) {
	return uc.adapterStorage.CreateContactIntoGroup(ctx, groupID, contacts...)
}

func (uc *UseCase) AddContactToGroup(ctx context.Context, groupID, contactID uuid.UUID) error {
	return uc.adapterStorage.AddContactsToGroup(ctx, groupID, contactID)
}

func (uc *UseCase) DeleteContactFromGroup(ctx context.Context, groupID, contactID uuid.UUID) error {
	return uc.adapterStorage.DeleteContactFromGroup(ctx, groupID, contactID)
}
