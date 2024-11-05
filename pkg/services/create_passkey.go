package services

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/a-novel/uservice-passkeys/pkg/dao"
)

var (
	ErrInvalidCreatePasskeyRequest = errors.New("invalid create passkey request")
	ErrCreatePasskey               = errors.New("create passkey")
)

var createPasskeyValidate = validator.New(validator.WithRequiredStructEnabled())

type CreatePasskeyRequest struct {
	Namespace string                 `validate:"required,min=1,max=256"`
	Passkey   string                 `validate:"required,min=4,max=4096"`
	Reward    map[string]interface{} `validate:"omitempty"`
	ExpiresIn *time.Duration         `validate:"omitempty"`
}

type CreatePasskeyResponse struct {
	ID        string
	Namespace string
	Reward    map[string]interface{}
	ExpiresAt *time.Time
	CreatedAt time.Time
}

type CreatePasskey interface {
	Exec(ctx context.Context, data *CreatePasskeyRequest) (*CreatePasskeyResponse, error)
}

type createPasskeyImpl struct {
	dao dao.CreatePasskey
}

func (service *createPasskeyImpl) Exec(
	ctx context.Context, data *CreatePasskeyRequest,
) (*CreatePasskeyResponse, error) {
	if err := createPasskeyValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidCreatePasskeyRequest, err)
	}

	request := &dao.CreatePasskeyRequest{
		Namespace: data.Namespace,
		Passkey:   data.Passkey,
		Reward:    data.Reward,
		ExpiresAt: ExpiresInToTime(data.ExpiresIn),
	}

	res, err := service.dao.Exec(ctx, uuid.New(), time.Now(), request)
	if err != nil {
		return nil, errors.Join(ErrCreatePasskey, err)
	}

	return &CreatePasskeyResponse{
		ID:        res.ID.String(),
		Namespace: res.Namespace,
		Reward:    res.Reward,
		ExpiresAt: res.ExpiresAt,
		CreatedAt: res.CreatedAt,
	}, nil
}

func NewCreatePasskey(dao dao.CreatePasskey) CreatePasskey {
	return &createPasskeyImpl{dao: dao}
}
