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

type UpdatePasskeyRequest struct {
	Namespace string
	Passkey   string
	Reward    map[string]interface{}
	ExpiresAt *time.Time
}

type UpdatePasskey interface {
	Exec(ctx context.Context, id uuid.UUID, now time.Time, request *UpdatePasskeyRequest) (*entities.Passkey, error)
}

type updatePasskeyImpl struct {
	database bun.IDB
}

func (dao *updatePasskeyImpl) Exec(
	ctx context.Context, passkeyID uuid.UUID, now time.Time, request *UpdatePasskeyRequest,
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
		UpdatedAt:    &now,
	}

	rows, err := dao.database.NewUpdate().
		Model(model).
		WherePK().
		ExcludeColumn("created_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("get rows affected: %w", err)
	}

	if affected == 0 {
		return nil, ErrPasskeyNotFound
	}

	return model, nil
}

func NewUpdatePasskey(database bun.IDB) UpdatePasskey {
	return &updatePasskeyImpl{database: database}
}
