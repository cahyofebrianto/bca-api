package context_test

import (
	"context"
	"testing"

	bcaCtx "github.com/purwaren/bca-api/context"
	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	ctx := bcaCtx.With(
		context.Background(),
		bcaCtx.HTTPReqID("httpReqID01"),
		bcaCtx.HTTPSessID("httpSessID02"),
		bcaCtx.BCASessID("bcaSessID03"))

	require.Equal(t, "httpReqID01", ctx.Value("httpReqID"))
	require.Equal(t, "httpSessID02", ctx.Value("httpSessID"))
	require.Equal(t, "bcaSessID03", ctx.Value("bcaSessID"))
}
