package bca

import (
	"context"
	"os"

	"github.com/avast/retry-go"
	"github.com/juju/errors"
	bcaCtx "github.com/purwaren/bca-api/context"
	"github.com/purwaren/bca-api/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type BCA struct {
	api    *api
	config Config
}

var MaxRetryAttempts int = 2

func New(config Config) *BCA {
	bca := BCA{
		config: config,
		api:    newAPI(config),
	}

	logger.SetOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {

		fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename: config.LogPath,
			MaxSize:  500, // megabytes
			// MaxBackups: 3,
			// MaxAge:     28, // days
		})
		stdoutWriteSyncer := zapcore.AddSync(os.Stdout)

		return zapcore.NewCore(
			zapcore.NewJSONEncoder(logger.DefaultEncoderConfig),
			zapcore.NewMultiWriteSyncer(fileWriteSyncer, stdoutWriteSyncer),
			zap.InfoLevel,
		)

		// return core
	}))

	return &bca
}

var ErrESB14009 = errors.New("Custom err. Meaning auth err from BCA API (ESB-14-009)")

func errorIfErrCodeESB14009(dtoError Error) error {
	if dtoError.ErrorCode == "ESB-14-009" {
		return ErrESB14009
	}
	return nil
}

func (b *BCA) retryDecision(ctx context.Context) func(err error) bool {
	return func(err error) bool {
		return err == ErrESB14009
	}
}

func (b *BCA) retryOptions(ctx context.Context) []retry.Option {
	return []retry.Option{
		retry.Attempts(2),
		retry.RetryIf(b.retryDecision(ctx)),
		retry.OnRetry(func(n uint, err error) {
			b.log(ctx).Infof("=== START ON RETRY === [Attempts: %d Err: %+v]", n, err)
			b.DoAuthentication(ctx)
			b.log(ctx).Infof("=== END ON RETRY ===")
		}),
	}
}

// === misc func ===

func (b *BCA) log(ctx context.Context) *zap.SugaredLogger {
	return logger.Logger(bcaCtx.With(ctx, bcaCtx.BCASessID(b.api.bcaSessID)))
}
