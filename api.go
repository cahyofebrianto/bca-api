package bca

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/juju/errors"
	"github.com/purwaren/bca-api/logger"
	"go.uber.org/zap"

	bcaCtx "github.com/purwaren/bca-api/context"
)

type API struct {
	config              Config
	httpClient          *http.Client // for postGetToken only
	retryablehttpClient *retryablehttp.Client

	mutex       sync.Mutex
	accessToken string
	bcaSessID   string
}

func NewAPI(config Config) *API {

	httpClient := cleanhttp.DefaultPooledClient()
	retryablehttpClient := retryablehttp.NewClient()

	api := API{config: config,
		httpClient:          httpClient,
		retryablehttpClient: retryablehttpClient,
	}

	return &api
}

func (api *API) SetAccessTokenAndSessID(accessToken, bcaSessID string) {
	api.mutex.Lock()
	defer api.mutex.Unlock()

	api.accessToken = accessToken
	api.bcaSessID = bcaSessID
}

func (api *API) PostGetToken(ctx context.Context) (*AuthToken, error) {
	urlTarget, err := buildURL(api.config.URL, "/api/oauth/token", url.Values{})
	if err != nil {
		return nil, errors.Trace(err)
	}

	form := url.Values{"grant_type": []string{"client_credentials"}}
	bodyReq := strings.NewReader(form.Encode())

	req, err := http.NewRequest(http.MethodPost, urlTarget, bodyReq)
	if err != nil {
		return nil, errors.Trace(err)
	}
	req = req.WithContext(ctx)

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(api.config.ClientID, api.config.ClientSecret)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()

	bodyRespBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}

	api.log(ctx).Info(resp.StatusCode)
	api.log(ctx).Info(string(bodyRespBytes))

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyRespBytes))

	var dtoResp AuthToken
	err = json.NewDecoder(resp.Body).Decode(&dtoResp)

	if err != nil {
		return nil, errors.Trace(err)
	}

	return &dtoResp, nil
}

// Generic POST request to API
func (api *API) Call(ctx context.Context, httpMethod string, path string, bodyReqPayload []byte, dtoResp interface{}) (err error) {
	// urlQuery := url.Values{"access_token": []string{api.accessToken}}
	urlQuery := url.Values{}
	urlTarget, err := buildURL(api.config.URL, path, urlQuery)
	if err != nil {
		return errors.Trace(err)
	}

	req, err := retryablehttp.NewRequest(httpMethod, urlTarget, bytes.NewBuffer(bodyReqPayload))
	if err != nil {
		return errors.Trace(err)
	}
	req = req.WithContext(ctx)

	req.Header.Set("content-type", "application/json")

	req.Header.Set("Authorization", "Bearer "+api.accessToken)
	req.Header.Set("Origin", api.config.OriginHost)
	req.Header.Set("X-BCA-Key", api.config.APIKey)

	timestamp := time.Now().Format("2006-01-02T15:04:05.999Z07:00")
	req.Header.Set("X-BCA-Timestamp", timestamp)

	signature := generateSignature(api.config.APISecret, httpMethod, path, api.accessToken, string(bodyReqPayload), timestamp)
	req.Header.Set("X-BCA-Signature", signature)

	resp, err := api.retryablehttpClient.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()

	bodyRespBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}

	api.log(ctx).Info(httpMethod + " " + urlTarget + " " + timestamp)
	api.log(ctx).Info(resp.StatusCode)
	api.log(ctx).Info(string(bodyRespBytes))

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyRespBytes))

	err = json.NewDecoder(resp.Body).Decode(&dtoResp)

	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

// === misc func ===
func (api *API) log(ctx context.Context) *zap.SugaredLogger {
	return logger.Logger(bcaCtx.WithBCASessID(ctx, api.bcaSessID))
}

func buildURL(baseUrl, paths string, query url.Values) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", errors.Trace(err)
	}

	u.Path = path.Join(u.Path, paths)
	u.RawQuery = query.Encode()

	return u.String(), nil
}
