package bca

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBCA_DoAuthentication_integration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	givenConfig := Config{
		URL:          os.Getenv("URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	bca := New(givenConfig)

	// resp based on sandbox resp
	dtoResp, err := bca.DoAuthentication(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, dtoResp.Error)
}
