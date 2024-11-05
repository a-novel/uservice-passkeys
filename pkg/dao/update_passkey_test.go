package dao_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	anoveldb "github.com/a-novel/golib/database"

	"github.com/a-novel/uservice-passkeys/migrations"
	"github.com/a-novel/uservice-passkeys/pkg/dao"
	"github.com/a-novel/uservice-passkeys/pkg/entities"
	"github.com/a-novel/uservice-passkeys/pkg/lib"
)

func TestUpdatePasskey(t *testing.T) {
	password1 := "password1"

	encryptedPassword1, err := lib.GenerateFromPassword(password1, lib.DefaultGenerateParams)
	require.NoError(t, err)

	fixtures := []interface{}{
		&entities.Passkey{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Namespace:    "namespace",
			EncryptedKey: encryptedPassword1,
			Reward:       map[string]interface{}{"key": "value"},
			ExpiresAt:    lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
			CreatedAt:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	testCases := []struct {
		name string

		id      uuid.UUID
		now     time.Time
		request *dao.UpdatePasskeyRequest

		expect    *entities.Passkey
		expectErr error
	}{
		{
			name: "Update",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),

			request: &dao.UpdatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
				Reward:    map[string]interface{}{"new-key": "new-value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 2, 0, 0, 0, 0, time.UTC)),
			},

			expect: &entities.Passkey{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Namespace: "namespace",
				Reward:    map[string]interface{}{"new-key": "new-value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 2, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "Update/RemoveUnnecessary",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),

			request: &dao.UpdatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
			},

			expect: &entities.Passkey{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Namespace: "namespace",
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "NotFound",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),

			request: &dao.UpdatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
			},

			expectErr: dao.ErrPasskeyNotFound,
		},
	}

	database, closer, err := anoveldb.OpenTestDB(&migrations.SQLMigrations)
	require.NoError(t, err)
	defer closer()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			transaction := anoveldb.BeginTestTX(database, fixtures)
			defer anoveldb.RollbackTestTX(transaction)

			updatePasskeyDAO := dao.NewUpdatePasskey(transaction)

			result, err := updatePasskeyDAO.Exec(context.Background(), testCase.id, testCase.now, testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)

			if testCase.expect == nil {
				require.Nil(t, result)
			} else {
				require.Equal(t, testCase.expect.ID, result.ID)
				require.Equal(t, testCase.expect.Namespace, result.Namespace)
				require.Equal(t, testCase.expect.Reward, result.Reward)
				require.Equal(t, testCase.expect.ExpiresAt, result.ExpiresAt)
				require.Equal(t, testCase.expect.CreatedAt, result.CreatedAt)
				require.Equal(t, testCase.expect.UpdatedAt, result.UpdatedAt)

				matching, err := lib.ComparePasswordAndHash(testCase.request.Passkey, result.EncryptedKey)
				require.NoError(t, err)
				require.True(t, matching)
			}
		})
	}
}
