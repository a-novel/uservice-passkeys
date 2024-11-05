package services_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/uservice-passkeys/pkg/dao"
	daomocks "github.com/a-novel/uservice-passkeys/pkg/dao/mocks"
	"github.com/a-novel/uservice-passkeys/pkg/entities"
	"github.com/a-novel/uservice-passkeys/pkg/services"
)

func TestCreatePasskey(t *testing.T) {
	testCases := []struct {
		name string

		request *services.CreatePasskeyRequest

		shouldCallCreatePasskeyDAO bool
		passkeyDAOResp             *entities.Passkey
		passkeyDAOErr              error

		expect    *services.CreatePasskeyResponse
		expectErr error
	}{
		{
			name: "OK",

			request: &services.CreatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresIn: lo.ToPtr(time.Hour * 24),
			},

			shouldCallCreatePasskeyDAO: true,
			passkeyDAOResp: &entities.Passkey{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Namespace: "namespace",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			},

			expect: &services.CreatePasskeyResponse{
				ID:        "00000000-0000-0000-0000-000000000002",
				Namespace: "namespace",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "OK/Minimal",

			request: &services.CreatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
			},

			shouldCallCreatePasskeyDAO: true,
			passkeyDAOResp: &entities.Passkey{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Namespace: "namespace",
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			},

			expect: &services.CreatePasskeyResponse{
				ID:        "00000000-0000-0000-0000-000000000002",
				Namespace: "namespace",
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "DAO/Error",

			request: &services.CreatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresIn: lo.ToPtr(time.Hour * 24),
			},

			shouldCallCreatePasskeyDAO: true,

			passkeyDAOErr: errors.New("uwups"),

			expectErr: services.ErrCreatePasskey,
		},
		{
			name: "InvalidRequest",

			request: &services.CreatePasskeyRequest{
				Passkey: "passkey",
			},

			expectErr: services.ErrInvalidCreatePasskeyRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			createPasskeyDAO := daomocks.NewMockCreatePasskey(t)

			if testCase.shouldCallCreatePasskeyDAO {
				createPasskeyDAO.
					On(
						"Exec",
						context.Background(),
						mock.MatchedBy(func(id uuid.UUID) bool { return id != uuid.Nil }),
						mock.MatchedBy(func(at time.Time) bool { return at.Unix() > 0 }),
						mock.MatchedBy(func(data *dao.CreatePasskeyRequest) bool {
							baseCHeck := data.Namespace == testCase.request.Namespace &&
								data.Passkey == testCase.request.Passkey &&
								reflect.DeepEqual(data.Reward, testCase.request.Reward)

							if testCase.request.ExpiresIn == nil {
								return baseCHeck && data.ExpiresAt == nil
							}

							return baseCHeck &&
								data.ExpiresAt != nil &&
								data.ExpiresAt.After(time.Now().Add(*testCase.request.ExpiresIn-time.Minute)) &&
								data.ExpiresAt.Before(time.Now().Add(*testCase.request.ExpiresIn+time.Minute))
						}),
					).
					Return(testCase.passkeyDAOResp, testCase.passkeyDAOErr)
			}

			service := services.NewCreatePasskey(createPasskeyDAO)
			response, err := service.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, response)

			createPasskeyDAO.AssertExpectations(t)
		})
	}
}
