package dao_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	anoveldb "github.com/a-novel/golib/database"
	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/formatters"

	"github.com/a-novel/uservice-passkeys/migrations"
	"github.com/a-novel/uservice-passkeys/pkg/dao"
	"github.com/a-novel/uservice-passkeys/pkg/entities"
	"github.com/a-novel/uservice-passkeys/pkg/lib"
)

func TestDeletePasskey(t *testing.T) {
	password1 := "password1"
	password2 := "password2"

	encryptedPassword1, err := lib.GenerateFromPassword(password1, lib.DefaultGenerateParams)
	require.NoError(t, err)
	encryptedPassword2, err := lib.GenerateFromPassword(password2, lib.DefaultGenerateParams)
	require.NoError(t, err)

	fixtures := []interface{}{
		&entities.Passkey{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Namespace:    "namespace",
			EncryptedKey: encryptedPassword1,
			Reward:       map[string]interface{}{"key": "value"},
			ExpiresAt:    lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
			CreatedAt:    time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:    lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
		},
		&entities.Passkey{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Namespace:    "namespace-2",
			EncryptedKey: encryptedPassword2,
			ExpiresAt:    lo.ToPtr(time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)),
			CreatedAt:    time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	testCases := []struct {
		name string

		request *dao.DeletePasskeyRequest

		expect    *entities.Passkey
		expectErr error
	}{
		{
			name: "Delete",

			request: &dao.DeletePasskeyRequest{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Namespace: "namespace",
			},

			expect: &entities.Passkey{
				ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Namespace:    "namespace",
				EncryptedKey: encryptedPassword1,
				Reward:       map[string]interface{}{"key": "value"},
				ExpiresAt:    lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt:    time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:    lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "Delete/Expired",

			request: &dao.DeletePasskeyRequest{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Namespace: "namespace-2",
			},

			expect: &entities.Passkey{
				ID:           uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Namespace:    "namespace-2",
				EncryptedKey: encryptedPassword2,
				ExpiresAt:    lo.ToPtr(time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt:    time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Delete/NotFound",

			request: &dao.DeletePasskeyRequest{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Namespace: "namespace",
			},

			expectErr: dao.ErrPasskeyNotFound,
		},
		{
			name: "Delete/WithPassword",

			request: &dao.DeletePasskeyRequest{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Namespace: "namespace",
				RawKey:    &password1,
			},

			expect: &entities.Passkey{
				ID:           uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Namespace:    "namespace",
				EncryptedKey: encryptedPassword1,
				Reward:       map[string]interface{}{"key": "value"},
				ExpiresAt:    lo.ToPtr(time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt:    time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:    lo.ToPtr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "Delete/WithPassword/BadPassword",

			request: &dao.DeletePasskeyRequest{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Namespace: "namespace",
				RawKey:    &password2,
			},

			expectErr: dao.ErrInvalidPasskey,
		},
	}

	database, closer, err := anoveldb.OpenTestDB(nil)
	require.NoError(t, err)
	defer closer()

	require.NoError(t, anoveldb.FreezeTime(database, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))

	formatter := formatters.NewConsoleFormatter(loggers.NewSTDOut(), true)
	require.NoError(t, anoveldb.Migrate(database, migrations.SQLMigrations, formatter))

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			transaction := anoveldb.BeginTestTX(database, fixtures)
			defer anoveldb.RollbackTestTX(transaction)

			deletePasskeyDAO := dao.NewDeletePasskey(transaction)

			result, err := deletePasskeyDAO.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, result)
		})
	}
}
