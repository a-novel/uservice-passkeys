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
	ErrInvalidGetPasskeyRequest = errors.New("invalid get passkey request")
	ErrGetPasskey               = errors.New("get passkey")
)

var getPasskeyValidate = validator.New(validator.WithRequiredStructEnabled())

type GetPasskeyRequest struct {
	ID        string `validate:"required,len=36"`
	Namespace string `validate:"required,min=1,max=256"`
	Passkey   string `validate:"required_if=Validate true,omitempty,min=4,max=4096"`
	Validate  bool   `validate:"omitempty"`
}

type GetPasskeyResponse struct {
	ID        string
	Namespace string
	Reward    map[string]interface{}
	ExpiresAt *time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type GetPasskey interface {
	Exec(ctx context.Context, data *GetPasskeyRequest) (*GetPasskeyResponse, error)
}

type getPasskeyImpl struct {
	dao dao.GetPasskey
}

func (service *getPasskeyImpl) Exec(
	ctx context.Context, data *GetPasskeyRequest,
) (*GetPasskeyResponse, error) {
	if err := getPasskeyValidate.Struct(data); err != nil {
		return nil, errors.Join(ErrInvalidGetPasskeyRequest, err)
	}

	passkeyID, err := uuid.Parse(data.ID)
	if err != nil {
		return nil, errors.Join(ErrInvalidGetPasskeyRequest, fmt.Errorf("uuid value: '%s': %w", data.ID, err))
	}

	request := &dao.GetPasskeyRequest{
		ID:        passkeyID,
		Namespace: data.Namespace,
		RawKey:    lo.Ternary[*string](data.Validate, &data.Passkey, nil),
	}

	res, err := service.dao.Exec(ctx, request)
	if err != nil {
		return nil, errors.Join(ErrGetPasskey, err)
	}

	return &GetPasskeyResponse{
		ID:        res.ID.String(),
		Namespace: res.Namespace,
		Reward:    res.Reward,
		ExpiresAt: res.ExpiresAt,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func NewGetPasskey(dao dao.GetPasskey) GetPasskey {
	return &getPasskeyImpl{dao: dao}
}
