package bca_test

import (
	"context"
	"os"
	"testing"

	"github.com/purwaren/bca-api"
	"github.com/stretchr/testify/require"
)

func TestBCA_Fire_integration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	t.Run("FireInquiryAccount", func(t *testing.T) {
		givenConfig := bca.Config{
			URL:          os.Getenv("URL"),
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),

			CorporateID: os.Getenv("CORPORATE_ID"),

			APIKey:    os.Getenv("API_KEY"),
			APISecret: os.Getenv("API_SECRET"),

			OriginHost: os.Getenv("ORIGIN_HOST"),
		}

		givenDtoReq := bca.InquiryAccountRequest{
			Authentication: bca.Authentication{
				CorporateID: "DUMMYI",
				AccessCode:  "Kw5oTuF12dseSH44Y8ww",
				BranchCode:  "BCA001",
				UserID:      "BCAUSERID",
				LocalID:     "40115"},
			BeneficiaryDetails: bca.InquiryAccountRequestBeneficiaryDetails{
				BankCodeType:  "BIC",
				BankCodeValue: "CENAIDJAXXX",
				AccountNumber: "0106666011"},
		}

		b := bca.New(givenConfig)
		// resp based on sandbox doc
		dtoResp, err := b.FireInquiryAccount(context.Background(), givenDtoReq)

		require.NoError(t, err)
		require.Empty(t, dtoResp.Error)
	})

}
