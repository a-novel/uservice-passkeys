package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/a-novel/uservice-passkeys/pkg/dao"
)

var (
	ErrInvalidUpdatePasskeyRequest = errors.New("invalid update passkey request")
	ErrUpdatePasskey               = errors.New("update passkey")
)

var updatePasskeyValidate = validator.New(validator.WithRequiredStructEnabled())

type UpdatePasskeyRequest struct {
	ID        string                 `validate:"required,len=36"`
	Namespace string                 `validate:"required,min=1,max=256"`
	Passkey   string                 `validate:"required,min=4,max=4096"`
	Reward    map[string]interface{} `validate:"omitempty"`
	ExpiresIn *time.Duration         `validate:"omitempty"`
}

type UpdatePasskeyResponse struct {
	ID        string
	Namespace string
	Reward    map[string]interface{}
	ExpiresAt *time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UpdatePasskey interface {
	Exec(ctx context.Context, data *UpdatePasskeyRequest) (*UpdatePasskeyResponse, error)
}

type updatePasskeyImpl struct {
	dao dao.UpdatePasskey
}

func (service *updatePasskeyImpl) Exec(
	ctx context.Context, data *UpdatePasskeyRequest,
) (*UpdatePasskeyResponse, error) {
	if err := updatePasskeyValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidUpdatePasskeyRequest, err)
	}

	passkeyID, err := uuid.Parse(data.ID)
	if err != nil {
		return nil, errors.Join(ErrInvalidUpdatePasskeyRequest, fmt.Errorf("uuid value: '%s': %w", data.ID, err))
	}

	request := &dao.UpdatePasskeyRequest{
		Namespace: data.Namespace,
		Passkey:   data.Passkey,
		Reward:    data.Reward,
		ExpiresAt: ExpiresInToTime(data.ExpiresIn),
	}

	res, err := service.dao.Exec(ctx, passkeyID, time.Now(), request)
	if err != nil {
		return nil, errors.Join(ErrUpdatePasskey, err)
	}

	return &UpdatePasskeyResponse{
		ID:        res.ID.String(),
		Namespace: res.Namespace,
		Reward:    res.Reward,
		ExpiresAt: res.ExpiresAt,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func NewUpdatePasskey(dao dao.UpdatePasskey) UpdatePasskey {
	return &updatePasskeyImpl{dao: dao}
}
