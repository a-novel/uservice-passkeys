package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/uservice-passkeys/pkg/dao"
	daomocks "github.com/a-novel/uservice-passkeys/pkg/dao/mocks"
	"github.com/a-novel/uservice-passkeys/pkg/entities"
	"github.com/a-novel/uservice-passkeys/pkg/services"
)

func TestGetPasskey(t *testing.T) {
	testCases := []struct {
		name string

		request *services.GetPasskeyRequest

		shouldCallDeletePasskeyDAO bool
		passkeyDAOResp             *entities.Passkey
		passkeyDAOErr              error

		expect    *services.GetPasskeyResponse
		expectErr error
	}{
		{
			name: "OK",

			request: &services.GetPasskeyRequest{
				ID:        "00000000-0000-0000-0000-000000000001",
				Namespace: "namespace",
			},

			shouldCallDeletePasskeyDAO: true,
			passkeyDAOResp: &entities.Passkey{
				ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Namespace:    "namespace",
				EncryptedKey: "encryptedKey",
				Reward:       map[string]interface{}{"key": "value"},
				ExpiresAt:    lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt:    time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:    lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},

			expect: &services.GetPasskeyResponse{
				ID:        "00000000-0000-0000-0000-000000000001",
				Namespace: "namespace",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "OK/WithPassword",

			request: &services.GetPasskeyRequest{
				ID:        "00000000-0000-0000-0000-000000000001",
				Namespace: "namespace",
				Passkey:   "passkey",
				Validate:  true,
			},

			shouldCallDeletePasskeyDAO: true,
			passkeyDAOResp: &entities.Passkey{
				ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Namespace:    "namespace",
				EncryptedKey: "encryptedKey",
				Reward:       map[string]interface{}{"key": "value"},
				ExpiresAt:    lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt:    time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:    lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},

			expect: &services.GetPasskeyResponse{
				ID:        "00000000-0000-0000-0000-000000000001",
				Namespace: "namespace",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "Error/WithPassword/NoPassword",

			request: &services.GetPasskeyRequest{
				ID:        "00000000-0000-0000-0000-000000000001",
				Namespace: "namespace",
				Validate:  true,
			},

			expectErr: services.ErrInvalidGetPasskeyRequest,
		},
		{
			name: "Error/InvalidID",

			request: &services.GetPasskeyRequest{
				ID:        "00000000x0000x0000x0000x000000000001",
				Namespace: "namespace",
			},

			expectErr: services.ErrInvalidGetPasskeyRequest,
		},
		{
			name: "DAO/Error",

			request: &services.GetPasskeyRequest{
				ID:        "00000000-0000-0000-0000-000000000001",
				Namespace: "namespace",
			},

			shouldCallDeletePasskeyDAO: true,
			passkeyDAOErr:              errors.New("uwups"),

			expectErr: services.ErrGetPasskey,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			getPasskeyDAO := daomocks.NewMockGetPasskey(t)

			if testCase.shouldCallDeletePasskeyDAO {
				getPasskeyDAO.
					On(
						"Exec",
						context.Background(),
						&dao.GetPasskeyRequest{
							ID:        uuid.MustParse(testCase.request.ID),
							Namespace: testCase.request.Namespace,
							RawKey:    lo.Ternary[*string](testCase.request.Validate, &testCase.request.Passkey, nil),
						},
					).
					Return(testCase.passkeyDAOResp, testCase.passkeyDAOErr)
			}

			service := services.NewGetPasskey(getPasskeyDAO)
			resp, err := service.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			getPasskeyDAO.AssertExpectations(t)
		})
	}
}
