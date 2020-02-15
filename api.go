package bca

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/juju/errors"
	"github.com/lithammer/shortuuid"
	"github.com/purwaren/bca-api/logger"
	"go.uber.org/zap"

	bcaCtx "github.com/purwaren/bca-api/context"
)

var (
	httpHeaderChannelID    string = "ChannelID"
	httpHeaderCredentialID string = "CredentialID"
)

type api struct {
	config     Config
	httpClient *http.Client // for postGetToken only

	mutex       sync.Mutex
	accessToken string
	bcaSessID   string
}

func newAPI(config Config) *api {

	httpClient := cleanhttp.DefaultPooledClient()

	api := api{config: config,
		httpClient: httpClient,
	}

	return &api
}

func (api *api) setAccessToken(accessToken string) {
	api.mutex.Lock()
	defer api.mutex.Unlock()

	newSessID := shortuuid.New()

	api.accessToken = accessToken
	api.bcaSessID = newSessID
}

// === AUTH ===

func (api *api) postGetToken(ctx context.Context) (*AuthToken, error) {
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

// === BANKING ===
func (api *api) bankingGetBalance(ctx context.Context, dtoReq BalanceInfoRequest) (*BalanceInfoResponse, error) {
	path := fmt.Sprintf("/banking/v3/corporates/%s/accounts/%s", api.config.CorporateID, dtoReq.AccountNumber)

	var balanceInfoResp BalanceInfoResponse
	if err := api.call(ctx, http.MethodGet, path, nil, []byte(""), &balanceInfoResp); err != nil {
		return nil, errors.Trace(err)
	}
	return &balanceInfoResp, nil
}

func (api *api) bankingPostFundTransfer(ctx context.Context, dtoReq FundTransferRequest) (*FundTransferResponse, error) {
	path := fmt.Sprintf("/banking/corporates/transfers")

	jsonReq, err := json.Marshal(dtoReq)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var fundTransferResp FundTransferResponse
	if err := api.call(ctx, http.MethodPost, path, nil, jsonReq, &fundTransferResp); err != nil {
		return nil, errors.Trace(err)
	}
	return &fundTransferResp, nil
}

func (api *api) bankingPostFundTransferDomestic(ctx context.Context, dtoReq FundTransferDomesticRequest) (*FundTransferDomesticResponse, error) {
	path := fmt.Sprintf("/banking/corporates/transfers/domestic")

	jsonReq, err := json.Marshal(dtoReq)
	if err != nil {
		return nil, errors.Trace(err)
	}

	headers := map[string]string{
		httpHeaderChannelID:    api.config.ChannelID,
		httpHeaderCredentialID: api.config.CredentialID,
	}

	var fundTransferDomesticResp FundTransferDomesticResponse
	if err := api.call(ctx, http.MethodPost, path, headers, jsonReq, &fundTransferDomesticResp); err != nil {
		return nil, errors.Trace(err)
	}
	return &fundTransferDomesticResp, nil
}

func (api *api) firePostInquiryAccount(ctx context.Context, dtoReq InquiryAccountRequest) (*InquiryAccountResponse, error) {
	path := fmt.Sprintf("/fire/accounts")

	jsonReq, err := json.Marshal(dtoReq)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var inquiryAccountResp InquiryAccountResponse
	if err := api.call(ctx, http.MethodPost, path, nil, jsonReq, &inquiryAccountResp); err != nil {
		return nil, errors.Trace(err)
	}
	return &inquiryAccountResp, nil
}

// Generic HTTP request to API
func (api *api) call(ctx context.Context, httpMethod string, path string, additionalHeader map[string]string, bodyReqPayload []byte, dtoResp interface{}) (err error) {
	// urlQuery := url.Values{"access_token": []string{api.accessToken}}
	urlQuery := url.Values{}
	urlTarget, err := buildURL(api.config.URL, path, urlQuery)
	if err != nil {
		return errors.Trace(err)
	}

	req, err := http.NewRequest(httpMethod, urlTarget, bytes.NewBuffer(bodyReqPayload))
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

	signature, _, err := GenerateSignature(api.config.APISecret, httpMethod, path, api.accessToken, string(bodyReqPayload), timestamp)
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Set("X-BCA-Signature", signature)

	api.log(ctx).Info(httpMethod + " " + urlTarget + " " + timestamp)
	// api.log(ctx).Info("StrToSIGN: " + strToSign)
	// api.log(ctx).Info("SIGN: " + signature)

	for key, val := range additionalHeader {
		req.Header.Set(key, val)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()

	bodyRespBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}

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
func (api *api) log(ctx context.Context) *zap.SugaredLogger {
	return logger.Logger(bcaCtx.With(ctx, bcaCtx.BCASessID(api.bcaSessID)))
}

func buildURL(baseURL, paths string, query url.Values) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", errors.Trace(err)
	}

	u.Path = path.Join(u.Path, paths)
	u.RawQuery = query.Encode()

	return u.String(), nil
}
