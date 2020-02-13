package bca_test

import (
	"context"
	"os"
	"testing"

	"github.com/purwaren/bca-api"
	"github.com/stretchr/testify/require"
)

func TestBCA_Auth_integration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("DoAuthentication", func(t *testing.T) {
		givenConfig := bca.Config{
			URL:          os.Getenv("URL"),
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		}

		b := bca.New(givenConfig)

		// resp based on sandbox resp
		dtoResp, err := b.DoAuthentication(context.Background())
		require.NoError(t, err)
		require.Empty(t, dtoResp.Error)
	})
}
