package bca

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBCA_Banking_integration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("BankingGetBalance", func(t *testing.T) {
		givenConfig := Config{
			URL:          os.Getenv("URL"),
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),

			CorporateID: os.Getenv("CORPORATE_ID"),

			APIKey:    os.Getenv("API_KEY"),
			APISecret: os.Getenv("API_SECRET"),

			OriginHost: os.Getenv("ORIGIN_HOST"),
		}

		bca := New(givenConfig)

		// resp based on sandbox doc
		givenDtoReq := BalanceInfoRequest{AccountNumber: "0201245680"}
		dtoResp, err := bca.BankingGetBalance(context.Background(), givenDtoReq)

		require.NoError(t, err)
		require.Empty(t, dtoResp.Error)
	})

	t.Run("BankingFundTransfer", func(t *testing.T) {
		givenConfig := Config{
			URL:          os.Getenv("URL"),
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),

			CorporateID: os.Getenv("CORPORATE_ID"),

			APIKey:    os.Getenv("API_KEY"),
			APISecret: os.Getenv("API_SECRET"),

			OriginHost: os.Getenv("ORIGIN_HOST"),
		}

		bca := New(givenConfig)

		// resp based on sandbox doc
		givenDtoReq := FundTransferRequest{
			SourceAccountNumber:      "0201245680",
			TransactionID:            "00000001",
			TransactionDate:          time.Now().Format("2006-01-02"),
			ReferenceID:              "12345/PO/2016",
			CurrencyCode:             "IDR",
			Amount:                   100000.00,
			BeneficiaryAccountNumber: "0201245681",
			Remark1:                  "Transfer Test",
			Remark2:                  "Online Transfer",
		}
		dtoResp, err := bca.BankingFundTransfer(context.Background(), givenDtoReq)

		require.NoError(t, err)
		require.Empty(t, dtoResp.Error)
	})

	t.Run("BankingFundTransferDomestic", func(t *testing.T) {
		givenConfig := Config{
			URL:          os.Getenv("URL"),
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),

			CorporateID: os.Getenv("CORPORATE_ID"),

			APIKey:    os.Getenv("API_KEY"),
			APISecret: os.Getenv("API_SECRET"),

			ChannelID:    os.Getenv("CHANNEL_ID"),
			CredentialID: os.Getenv("CREDENTIAL_ID"),

			OriginHost: os.Getenv("ORIGIN_HOST"),
		}

		bca := New(givenConfig)

		// resp based on sandbox doc
		givenDtoReq := FundTransferDomesticRequest{
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
		dtoResp, err := bca.BankingFundTransferDomestic(context.Background(), givenDtoReq)

		require.NoError(t, err)
		require.Empty(t, dtoResp.Error)
	})
}
