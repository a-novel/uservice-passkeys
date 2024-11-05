package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel/uservice-passkeys/pkg/dao"
)

var (
	ErrInvalidDeletePasskeyRequest = errors.New("invalid delete passkey request")
	ErrDeletePasskey               = errors.New("delete passkey")
)

var deletePasskeyValidate = validator.New(validator.WithRequiredStructEnabled())

type DeletePasskeyRequest struct {
	ID        string `validate:"required,len=36"`
	Namespace string `validate:"required,min=1,max=256"`
	Passkey   string `validate:"required_if=Validate true,omitempty,min=4,max=4096"`
	Validate  bool   `validate:"omitempty"`
}

type DeletePasskeyResponse struct {
	ID        string
	Namespace string
	Reward    map[string]interface{}
	ExpiresAt *time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type DeletePasskey interface {
	Exec(ctx context.Context, data *DeletePasskeyRequest) (*DeletePasskeyResponse, error)
}

type deletePasskeyImpl struct {
	dao dao.DeletePasskey
}

func (service *deletePasskeyImpl) Exec(
	ctx context.Context, data *DeletePasskeyRequest,
) (*DeletePasskeyResponse, error) {
	if err := deletePasskeyValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidDeletePasskeyRequest, err)
	}

	passkeyID, err := uuid.Parse(data.ID)
	if err != nil {
		return nil, errors.Join(ErrInvalidDeletePasskeyRequest, fmt.Errorf("uuid value: '%s': %w", data.ID, err))
	}

	request := &dao.DeletePasskeyRequest{
		ID:        passkeyID,
		Namespace: data.Namespace,
		RawKey:    lo.Ternary[*string](data.Validate, &data.Passkey, nil),
	}

	res, err := service.dao.Exec(ctx, request)
	if err != nil {
		return nil, errors.Join(ErrDeletePasskey, err)
	}

	return &DeletePasskeyResponse{
		ID:        res.ID.String(),
		Namespace: res.Namespace,
		Reward:    res.Reward,
		ExpiresAt: res.ExpiresAt,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func NewDeletePasskey(dao dao.DeletePasskey) DeletePasskey {
	return &deletePasskeyImpl{dao: dao}
}
