package context

import "context"

const (
	HTTPReqIDKey  = "httpReqID"
	HTTPSessIDKey = "httpSessID"
	BCASessIDKey  = "bcaSessID"
)

type buildCtxFunc func(ctx context.Context) context.Context

func With(ctx context.Context, buildFuncs ...buildCtxFunc) context.Context {
	for _, buildCtxFunc := range buildFuncs {
		ctx = buildCtxFunc(ctx)
	}
	return ctx
}

func HTTPReqID(reqID string) buildCtxFunc {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, HTTPReqIDKey, reqID)
	}
}

func HTTPSessID(reqID string) buildCtxFunc {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, HTTPSessIDKey, reqID)
	}
}

func BCASessID(reqID string) buildCtxFunc {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, BCASessIDKey, reqID)
	}
}
