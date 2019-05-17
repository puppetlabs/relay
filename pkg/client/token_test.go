package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	token := Token("abc123")

	require.Equal(t, token.Bearer(), "Bearer: abc123")
}
