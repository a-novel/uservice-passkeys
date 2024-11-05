package lib_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/uservice-passkeys/pkg/lib"
)

func TestRandomString(t *testing.T) {
	str, err := lib.Random(10)
	require.NoError(t, err)
	require.Len(t, str, 10)

	str2, err := lib.Random(10)
	require.NoError(t, err)
	require.Len(t, str2, 10)

	require.NotEqual(t, str, str2)

	str3, err := lib.Random(32)
	require.NoError(t, err)
	require.Len(t, str3, 32)
}
