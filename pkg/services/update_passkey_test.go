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

func TestUpdatePasskey(t *testing.T) {
	testCases := []struct {
		name string

		request *services.UpdatePasskeyRequest

		shouldCallUpdatePasskeyDAO bool
		passkeyDAOResp             *entities.Passkey
		passkeyDAOErr              error

		expect    *services.UpdatePasskeyResponse
		expectErr error
	}{
		{
			name: "OK",

			request: &services.UpdatePasskeyRequest{
				ID:        "00000000-0000-0000-0000-000000000002",
				Namespace: "namespace",
				Passkey:   "passkey",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresIn: lo.ToPtr(time.Hour * 24),
			},

			shouldCallUpdatePasskeyDAO: true,
			passkeyDAOResp: &entities.Passkey{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Namespace: "namespace",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},

			expect: &services.UpdatePasskeyResponse{
				ID:        "00000000-0000-0000-0000-000000000002",
				Namespace: "namespace",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "OK/Minimal",

			request: &services.UpdatePasskeyRequest{
				ID:        "00000000-0000-0000-0000-000000000002",
				Namespace: "namespace",
				Passkey:   "passkey",
			},

			shouldCallUpdatePasskeyDAO: true,
			passkeyDAOResp: &entities.Passkey{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Namespace: "namespace",
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},

			expect: &services.UpdatePasskeyResponse{
				ID:        "00000000-0000-0000-0000-000000000002",
				Namespace: "namespace",
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "DAO/Error",

			request: &services.UpdatePasskeyRequest{
				ID:        "00000000-0000-0000-0000-000000000002",
				Namespace: "namespace",
				Passkey:   "passkey",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresIn: lo.ToPtr(time.Hour * 24),
			},

			shouldCallUpdatePasskeyDAO: true,

			passkeyDAOErr: errors.New("uwups"),

			expectErr: services.ErrUpdatePasskey,
		},
		{
			name: "InvalidRequest",

			request: &services.UpdatePasskeyRequest{
				Passkey: "passkey",
			},

			expectErr: services.ErrInvalidUpdatePasskeyRequest,
		},
		{
			name: "InvalidID",

			request: &services.UpdatePasskeyRequest{
				ID:        "00000000x0000x0000x0000x000000000002",
				Namespace: "namespace",
				Passkey:   "passkey",
			},

			expectErr: services.ErrInvalidUpdatePasskeyRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			updatePasskeyDAO := daomocks.NewMockUpdatePasskey(t)

			if testCase.shouldCallUpdatePasskeyDAO {
				updatePasskeyDAO.
					On(
						"Exec",
						context.Background(),
						uuid.MustParse(testCase.request.ID),
						mock.MatchedBy(func(at time.Time) bool { return at.Unix() > 0 }),
						mock.MatchedBy(func(data *dao.UpdatePasskeyRequest) bool {
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

			service := services.NewUpdatePasskey(updatePasskeyDAO)
			resp, err := service.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			updatePasskeyDAO.AssertExpectations(t)
		})
	}
}
