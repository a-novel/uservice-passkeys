package dao

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/uservice-passkeys/pkg/entities"
	"github.com/a-novel/uservice-passkeys/pkg/lib"
)

type DeletePasskeyRequest struct {
	ID        uuid.UUID
	Namespace string
	RawKey    *string
}

type DeletePasskey interface {
	Exec(ctx context.Context, request *DeletePasskeyRequest) (*entities.Passkey, error)
}

type deletePasskeyImpl struct {
	database bun.IDB
}

func (dao *deletePasskeyImpl) Exec(ctx context.Context, request *DeletePasskeyRequest) (*entities.Passkey, error) {
	model := &entities.Passkey{
		ID:        request.ID,
		Namespace: request.Namespace,
	}

	txErr := dao.database.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		rows, err := tx.NewDelete().
			Model(model).
			WherePK().
			Returning("*").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("exec query: %w", err)
		}

		affected, err := rows.RowsAffected()
		if err != nil {
			return fmt.Errorf("get rows affected: %w", err)
		}

		if affected == 0 {
			return ErrPasskeyNotFound
		}

		if request.RawKey != nil {
			match, err := lib.ComparePasswordAndHash(*request.RawKey, model.EncryptedKey)
			if err != nil {
				return fmt.Errorf("compare passkey: %w", err)
			}

			if !match {
				return ErrInvalidPasskey
			}
		}

		return nil
	})
	if txErr != nil {
		return nil, fmt.Errorf("exec transaction: %w", txErr)
	}

	return model, nil
}

func NewDeletePasskey(database bun.IDB) DeletePasskey {
	return &deletePasskeyImpl{database: database}
}
