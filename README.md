# Unofficial Go SDK for Bank BCA API

![build](https://github.com/purwaren/bca-api/workflows/Go/badge.svg?branch=master)

## (Currently) Supported API

- `POST /api/oauth/token` (`DoAuthentication`)
- `GET /banking/v3/corporates/<CorporateID>/accounts/<AccountNum>` (`BankingGetBalance`)
- `POST /banking/corporates/transfers` (`BankingFundTransfer`)
- `POST /banking/corporates/transfers/domestic` (`BankingFundTransferDomestic`)
- `POST /fire/accounts` (`FireInquiryAccount`)

For the detail, see [official documentation of BCA API](https://developer.bca.co.id/documentation/)

## Usage

NOTE: You don't have to explicitly do authentication before calling API. If got an auth error `Unauthorized` (`ErrorCode:ESB-14-009`), it will automatically retry failed API operation but `DoAuthentication` beforehand. Default max retry attempts is only 2.

```go
package main

import (
	"context"

	"github.com/lithammer/shortuuid"
	"github.com/purwaren/bca-api"
	bcaCtx "github.com/purwaren/bca-api/context"
)

func main() {
	cfg := bca.Config{
		URL: "https://sandbox.bca.co.id",
		ClientID:     "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		ClientSecret: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		APIKey:       "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		APISecret:    "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		CorporateID:  "BCAAPI2016",
		OriginHost:   "localhost",

		ChannelID:    "95051",
		CredentialID: "BCAAPI",

		LogPath: "bca.log",
	}

	api := bca.New(cfg)

	ctx := context.Background()
	ctx = bcaCtx.WithHTTPReqID(ctx, shortuuid.New())

	// bca.DoAuthentication(ctx) // <- You don't have to do this explicitly

	balanceInfoReq := bca.BalanceInfoRequest{AccountNumber: "0201245680"}
	api.BankingGetBalance(ctx, balanceInfoReq)

	fundTransferReq := bca.FundTransferRequest{
		SourceAccountNumber:      "0201245680",
		TransactionID:            "00000001",
		TransactionDate:          "2020-01-30",
		ReferenceID:              "12345/PO/2016",
		CurrencyCode:             "IDR",
		Amount:                   100000.00,
		BeneficiaryAccountNumber: "0201245681",
		Remark1:                  "Transfer Test",
		Remark2:                  "Online Transfer",
	}
	api.BankingFundTransfer(ctx, fundTransferReq)

	fundTransferDomesticReq := bca.FundTransferDomesticRequest{
		TransactionID:            "00000001",
		TransactionDate:          "2018-05-03",
		ReferenceID:              "12345/PO/2016",
		SourceAccountNumber:      "0201245680",
		BeneficiaryAccountNumber: "0201245501",
		BeneficiaryBankCode:      "BRONINJA",
		BeneficiaryName:          "Tester",
		Amount:                   100000.00,
		TransferType:             "LLG",
		BeneficiaryCustType:      "1",
		BeneficiaryCustResidence: "1",
		CurrencyCode:             "IDR",
		Remark1:                  "Transfer Test",
		Remark2:                  "Online Transfer",
	}
	api.BankingFundTransferDomestic(ctx, fundTransferDomesticReq)
}
```

## Contributing

Read the [Contribution Guide](CONTRIBUTING.md).
