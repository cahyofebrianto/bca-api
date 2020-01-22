package bca

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/lithammer/shortuuid"
	bcaCtx "github.com/purwaren/bca-api/context"
	"github.com/purwaren/bca-api/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type BCA struct {
	api    *api
	config Config

	mutex     sync.Mutex
	bcaSessID string
}

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

func (b *BCA) setAccessToken(accessToken string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	newSessID := shortuuid.New()
	b.api.setAccessTokenAndSessID(accessToken, newSessID)
	b.bcaSessID = newSessID
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

// === misc func ===

func (b *BCA) log(ctx context.Context) *zap.SugaredLogger {
	return logger.Logger(bcaCtx.WithBCASessID(ctx, b.bcaSessID))
}
