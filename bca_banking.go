package bca

import (
	"context"

	"github.com/avast/retry-go"
	"github.com/juju/errors"
	bcaCtx "github.com/purwaren/bca-api/context"
)

func (b *BCA) BankingGetBalance(ctx context.Context, dtoReq BalanceInfoRequest) (dtoResp *BalanceInfoResponse, err error) {
	ctx = bcaCtx.With(ctx, bcaCtx.BCASessID(b.api.bcaSessID))

	b.log(ctx).Info("=== START BANKING GET_BALANCE ===")
	b.log(ctx).Infof("REQUEST: %+v", dtoReq)

	retryOpts := b.retryOptions(ctx)
	err = retry.Do(func() error {
		if dtoResp, err = b.api.bankingGetBalance(ctx, dtoReq); err != nil {
			return err
		}
		return errorIfErrCodeESB14009(dtoResp.Error)
	}, retryOpts...)

	if err != nil {
		b.log(ctx).Error(errors.Details(err))
		return nil, errors.Trace(err)
	}

	b.log(ctx).Infof("RESPONSE: %+v", dtoResp)
	b.log(ctx).Info("=== END BANKING GET_BALANCE ===")

	return dtoResp, nil
}

func (b *BCA) BankingFundTransfer(ctx context.Context, dtoReq FundTransferRequest) (dtoResp *FundTransferResponse, err error) {
	ctx = bcaCtx.With(ctx, bcaCtx.BCASessID(b.api.bcaSessID))

	dtoReq.CorporateID = b.config.CorporateID

	b.log(ctx).Info("=== START BANKING FUND_TRANSFER ===")
	b.log(ctx).Infof("REQUEST: %+v", dtoReq)

	retryOpts := b.retryOptions(ctx)
	err = retry.Do(func() error {
		if dtoResp, err = b.api.bankingPostFundTransfer(ctx, dtoReq); err != nil {
			return err
		}
		return errorIfErrCodeESB14009(dtoResp.Error)
	}, retryOpts...)

	if err != nil {
		b.log(ctx).Error(errors.Details(err))
		return nil, errors.Trace(err)
	}

	b.log(ctx).Infof("RESPONSE: %+v", dtoResp)
	b.log(ctx).Info("=== END BANKING FUND_TRANSFER ===")

	return dtoResp, nil
}

func (b *BCA) BankingFundTransferDomestic(ctx context.Context, dtoReq FundTransferDomesticRequest) (dtoResp *FundTransferDomesticResponse, err error) {
	ctx = bcaCtx.With(ctx, bcaCtx.BCASessID(b.api.bcaSessID))

	b.log(ctx).Info("=== START BANKING FUND_TRANSFER_DOMESTIC ===")
	b.log(ctx).Infof("REQUEST: %+v", dtoReq)

	retryOpts := b.retryOptions(ctx)
	err = retry.Do(func() error {
		if dtoResp, err = b.api.bankingPostFundTransferDomestic(ctx, dtoReq); err != nil {
			return err
		}
		return errorIfErrCodeESB14009(dtoResp.Error)
	}, retryOpts...)

	if err != nil {
		b.log(ctx).Error(errors.Details(err))
		return nil, errors.Trace(err)
	}

	b.log(ctx).Infof("RESPONSE: %+v", dtoResp)
	b.log(ctx).Info("=== END BANKING FUND_TRANSFER_DOMESTIC ===")

	return dtoResp, nil
}

func (b *BCA) GenerateSignature(apiSecret, method, path, accessToken, requestBody, timestamp string) (signature string, strToSign string, err error) {
	return generateSignature(apiSecret, method, path, accessToken, requestBody, timestamp)
}
