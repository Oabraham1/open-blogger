package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	_, _, err := Connect()
	require.NoError(t, err)
}
