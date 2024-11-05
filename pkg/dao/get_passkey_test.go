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

func TestGetPasskey(t *testing.T) {
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
		},
		// Expired
		&entities.Passkey{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Namespace:    "namespace-2",
			EncryptedKey: encryptedPassword2,
			ExpiresAt:    lo.ToPtr(time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)),
			CreatedAt:    time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		// No expiry
		&entities.Passkey{
			ID:           uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Namespace:    "namespace-2",
			EncryptedKey: encryptedPassword2,
			Reward:       map[string]interface{}{"key": "value"},
			CreatedAt:    time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	testCases := []struct {
		name string

		request *dao.GetPasskeyRequest

		expect    *entities.Passkey
		expectErr error
	}{
		{
			name: "Get",

			request: &dao.GetPasskeyRequest{
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
			},
		},
		{
			name: "Get/NoExpiry",

			request: &dao.GetPasskeyRequest{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Namespace: "namespace-2",
			},

			expect: &entities.Passkey{
				ID:           uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				Namespace:    "namespace-2",
				EncryptedKey: encryptedPassword2,
				Reward:       map[string]interface{}{"key": "value"},
				CreatedAt:    time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Get/Expired",

			request: &dao.GetPasskeyRequest{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Namespace: "namespace-2",
			},

			expectErr: dao.ErrPasskeyNotFound,
		},
		{
			name: "NotFound",

			request: &dao.GetPasskeyRequest{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				Namespace: "namespace-2",
			},

			expectErr: dao.ErrPasskeyNotFound,
		},
		{
			name: "Get/WithPassword",

			request: &dao.GetPasskeyRequest{
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
			},
		},
		{
			name: "Get/WithPassword/BadPassword",

			request: &dao.GetPasskeyRequest{
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

	transaction := anoveldb.BeginTestTX(database, fixtures)
	defer anoveldb.RollbackTestTX(transaction)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			getPasskeyDAO := dao.NewGetPasskey(transaction)

			result, err := getPasskeyDAO.Exec(context.Background(), testCase.request)

			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, result)
		})
	}
}
