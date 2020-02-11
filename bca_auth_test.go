package bca

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBCA_Auth_integration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("DoAuthentication", func(t *testing.T) {
		givenConfig := Config{
			URL:          os.Getenv("URL"),
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		}

		bca := New(givenConfig)

		// resp based on sandbox resp
		dtoResp, err := bca.DoAuthentication(context.Background())
		require.NoError(t, err)
		require.Empty(t, dtoResp.Error)
	})
}
