package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/uservice-passkeys/pkg/entities"
	"github.com/a-novel/uservice-passkeys/pkg/lib"
)

type GetPasskeyRequest struct {
	ID        uuid.UUID
	Namespace string
	RawKey    *string
}

type GetPasskey interface {
	Exec(ctx context.Context, request *GetPasskeyRequest) (*entities.Passkey, error)
}

type getPasskeyImpl struct {
	database bun.IDB
}

func (dao *getPasskeyImpl) Exec(ctx context.Context, request *GetPasskeyRequest) (*entities.Passkey, error) {
	model := &entities.Passkey{
		ID:        request.ID,
		Namespace: request.Namespace,
	}

	err := dao.database.NewSelect().
		Model(model).
		WherePK().
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPasskeyNotFound
		}

		return nil, fmt.Errorf("exec query: %w", err)
	}

	if request.RawKey != nil {
		match, err := lib.ComparePasswordAndHash(*request.RawKey, model.EncryptedKey)
		if err != nil {
			return nil, fmt.Errorf("compare passkey: %w", err)
		}

		if !match {
			return nil, ErrInvalidPasskey
		}
	}

	return model, nil
}

func NewGetPasskey(database bun.IDB) GetPasskey {
	return &getPasskeyImpl{database: database}
}
