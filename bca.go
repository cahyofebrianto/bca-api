package bca

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/avast/retry-go"
	"github.com/hashicorp/go-retryablehttp"
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

	retryablehttpClient := retryablehttp.NewClient()
	retryablehttpClient.RetryMax = 1
	retryablehttpClient.CheckRetry = bca.retryPolicy

	bca.api.retryablehttpClient = retryablehttpClient

	return &bca
}

func (b *BCA) retryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		b.log(ctx).Infof("[Not Retry] Got error in context: %+v", ctx.Err())
		return false, ctx.Err()
	}

	bodyRespBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.log(ctx).Infof("[Not Retry] Failed to read response body: %+v", err)
		return false, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyRespBytes))

	var dtoResp Error
	if err := json.NewDecoder(bytes.NewReader(bodyRespBytes)).Decode(&dtoResp); err != nil {
		b.log(ctx).Infof("[Not Retry] Failed to decode error: %+v", err)
		return false, err
	}

	if resp.StatusCode == http.StatusUnauthorized || dtoResp.ErrorCode == "ESB-14-009" {
		b.log(ctx).Infof("[Retry] to auth. Got resp (code: %d, status: %s). Prev err: %+v. ErrorResp: %+v", resp.StatusCode, resp.Status, err, dtoResp)
		_, errAuth := b.DoAuthentication(ctx)

		return true, errAuth
	}

	return false, nil
}

// ========== RETRY V2 (FAIL as not elegant) ============================

func (b *BCA) callWithRetry(ctx context.Context, httpMethod string, path string, additionalHeader map[string]string, bodyReqPayload []byte, dtoResp interface{}) (err error) {
	attempts := 1

	for {
		attempts++
		if err = b.api.call(ctx, httpMethod, path, additionalHeader, bodyReqPayload, dtoResp); err != nil {
			return errors.Trace(err)
		}

		// check error
		dtoError, ok := dtoResp.(Error)
		if !ok {
			return
		}

		if dtoError.ErrorCode == "ESB-14-009" {
			b.log(ctx).Infof("[Retry] to auth")
			if _, err = b.DoAuthentication(ctx); err != nil {
				return errors.Trace(err)
			}
		}

		// retry decision
		if attempts >= MaxRetryAttempts {
			return
		}
	}

}

// ========== RETRY V3 (SUCCESS) ============================

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
