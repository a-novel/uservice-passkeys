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

func TestCreatePasskey(t *testing.T) {
	testCases := []struct {
		name string

		id      uuid.UUID
		now     time.Time
		request *dao.CreatePasskeyRequest

		expect    *entities.Passkey
		expectErr error
	}{
		{
			name: "Create",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),

			request: &dao.CreatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
			},

			expect: &entities.Passkey{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Namespace: "namespace",
				Reward:    map[string]interface{}{"key": "value"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Create/Minimal",

			id:  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			now: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),

			request: &dao.CreatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
			},

			expect: &entities.Passkey{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Namespace: "namespace",
				CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	database, closer, err := anoveldb.OpenTestDB(&migrations.SQLMigrations)
	require.NoError(t, err)
	defer closer()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			transaction := anoveldb.BeginTestTX[interface{}](database, nil)
			defer anoveldb.RollbackTestTX(transaction)

			createPasskeyDAO := dao.NewCreatePasskey(transaction)

			result, err := createPasskeyDAO.Exec(context.Background(), testCase.id, testCase.now, testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)

			if testCase.expect == nil {
				require.Nil(t, result)
			} else {
				require.Equal(t, testCase.expect.ID, result.ID)
				require.Equal(t, testCase.expect.Namespace, result.Namespace)
				require.Equal(t, testCase.expect.Reward, result.Reward)
				require.Equal(t, testCase.expect.ExpiresAt, result.ExpiresAt)
				require.Equal(t, testCase.expect.CreatedAt, result.CreatedAt)

				matching, err := lib.ComparePasswordAndHash(testCase.request.Passkey, result.EncryptedKey)
				require.NoError(t, err)
				require.True(t, matching)
			}
		})
	}
}
