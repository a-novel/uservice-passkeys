package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/uservice-passkeys/pkg/entities"
	"github.com/a-novel/uservice-passkeys/pkg/lib"
)

type CreatePasskeyRequest struct {
	Namespace string
	Passkey   string
	Reward    map[string]interface{}
	ExpiresAt *time.Time
}

type CreatePasskey interface {
	Exec(ctx context.Context, id uuid.UUID, now time.Time, request *CreatePasskeyRequest) (*entities.Passkey, error)
}

type createPasskeyImpl struct {
	database bun.IDB
}

func (dao *createPasskeyImpl) Exec(
	ctx context.Context, passkeyID uuid.UUID, now time.Time, request *CreatePasskeyRequest,
) (*entities.Passkey, error) {
	encrypted, err := lib.GenerateFromPassword(request.Passkey, lib.DefaultGenerateParams)
	if err != nil {
		return nil, fmt.Errorf("encrypt passkey: %w", err)
	}

	model := &entities.Passkey{
		ID:           passkeyID,
		Namespace:    request.Namespace,
		EncryptedKey: encrypted,
		Reward:       request.Reward,
		ExpiresAt:    request.ExpiresAt,
		CreatedAt:    now,
	}

	_, err = dao.database.NewInsert().Model(model).Returning("*").Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	return model, nil
}

func NewCreatePasskey(database bun.IDB) CreatePasskey {
	return &createPasskeyImpl{database: database}
}
