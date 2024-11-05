package lib_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/uservice-passkeys/pkg/lib"
)

func TestPassword(t *testing.T) {
	password := "password"

	encrypted, err := lib.GenerateFromPassword(password, lib.DefaultGenerateParams)
	require.NoError(t, err)
	require.NotEmpty(t, encrypted)

	testCases := []struct {
		name string

		password  string
		encrypted string

		expect    bool
		expectErr error
	}{
		{
			name: "OK",

			password:  password,
			encrypted: encrypted,

			expect: true,
		},
		{
			name: "WrongPassword",

			password:  "wrongpassword",
			encrypted: encrypted,

			expect: false,
		},
		{
			name: "Malformed/NotEnoughParts",

			password:  "password",
			encrypted: "malformed$",

			expectErr: lib.ErrInvalidHash,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ok, err := lib.ComparePasswordAndHash(testCase.password, testCase.encrypted)
			require.ErrorIs(t, testCase.expectErr, err)
			require.Equal(t, testCase.expect, ok)
		})
	}
}
