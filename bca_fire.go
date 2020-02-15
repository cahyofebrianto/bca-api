package bca

import (
	"context"

	"github.com/avast/retry-go"
	"github.com/juju/errors"
	bcaCtx "github.com/purwaren/bca-api/context"
)

// FireInquiryAccount inquiry BCA’s Account name or Other Bank Switching’s Account
func (b *BCA) FireInquiryAccount(ctx context.Context, dtoReq InquiryAccountRequest) (dtoResp *InquiryAccountResponse, err error) {
	ctx = bcaCtx.With(ctx, bcaCtx.BCASessID(b.api.bcaSessID))

	b.log(ctx).Info("=== START FIRE INQUIRY_ACCOUNT ===")
	b.log(ctx).Infof("REQUEST: %+v", dtoReq)

	retryOpts := b.retryOptions(ctx)
	err = retry.Do(func() error {
		if dtoResp, err = b.api.firePostInquiryAccount(ctx, dtoReq); err != nil {
			return err
		}
		return errorIfErrCodeESB14009(dtoResp.Error)
	}, retryOpts...)

	if err != nil {
		b.log(ctx).Error(errors.Details(err))
		return nil, errors.Trace(err)
	}

	b.log(ctx).Infof("RESPONSE: %+v", dtoResp)
	b.log(ctx).Info("=== END FIRE INQUIRY_ACCOUNT ===")

	return dtoResp, nil
}
