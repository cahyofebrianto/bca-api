package bca

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generateSignature(t *testing.T) {
	// Based on Authentication > Signature section in https://developer.bca.co.id/documentation/#signature

	ApiSecret := "22a2d25e-765d-41e1-8d29-da68dcb5698b"
	AccessToken := "lIWOt2p29grUo59bedBUrBY3pnzqQX544LzYPohcGHOuwn8AUEdUKS"

	type args struct {
		apiSecret   string
		method      string
		path        string
		accessToken string
		requestBody string
		timestamp   string
	}
	tests := []struct {
		name          string
		args          args
		wantSign      string
		wantStrToSign string
		wantErr       bool
	}{
		{name: "POST Example", args: args{
			apiSecret:   ApiSecret,
			method:      http.MethodPost,
			path:        "/banking/corporates/transfers",
			accessToken: AccessToken,
			requestBody: `
			{ 
				"CorporateID" : "BCAAPI2016",
					"SourceAccountNumber" : "0201245680",
					"TransactionID" : "00000001",
					"TransactionDate" : "2016-01-30",
					"ReferenceID" : "12345/PO/2016",
					"CurrencyCode" : "IDR",
					"Amount" : "100000.00",
					"BeneficiaryAccountNumber" : "0201245681",
					"Remark1" : "Transfer Test",
					"Remark2" : "Online Transfer"
				}
			`,
			timestamp: "2016-02-03T10:00:00.000+07:00",
		},
			wantSign:      "69ad66589ade078a30922a0848725cf153aecfcca82eba94e3270285b4a9c604",
			wantStrToSign: "POST:/banking/corporates/transfers:lIWOt2p29grUo59bedBUrBY3pnzqQX544LzYPohcGHOuwn8AUEdUKS:e3cf5797ac4ac02f7dad89ed2c5f5615c9884b2d802a504e4aebb76f45b8bdfb:2016-02-03T10:00:00.000+07:00",
			wantErr:       false,
		},
		{name: "GET Example", args: args{
			apiSecret:   ApiSecret,
			method:      http.MethodGet,
			path:        "/banking/v2/corporates/BCAAPI2016/accounts/0201245680/statements?StartDate=2016-09-01&EndDate=2016-09-01",
			accessToken: AccessToken,
			requestBody: "",
			timestamp:   "2016-02-03T10:00:00.000+07:00",
		},
			wantSign:      "3ac124303746d222387d4398dddf33201a384aa22137aa08f4d9843c6f467a48",
			wantStrToSign: "GET:/banking/v2/corporates/BCAAPI2016/accounts/0201245680/statements?EndDate=2016-09-01&StartDate=2016-09-01:lIWOt2p29grUo59bedBUrBY3pnzqQX544LzYPohcGHOuwn8AUEdUKS:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855:2016-02-03T10:00:00.000+07:00",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSign, gotStrToSign, err := generateSignature(tt.args.apiSecret, tt.args.method, tt.args.path, tt.args.accessToken, tt.args.requestBody, tt.args.timestamp)
			require.NoError(t, err)
			require.Equal(t, tt.wantStrToSign, gotStrToSign)
			require.Equal(t, tt.wantSign, gotSign)
		})
	}
}

func Test_canonicalize(t *testing.T) {
	// Based on Authentication > Signature section in https://developer.bca.co.id/documentation/#signature

	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Canonicalization Example", args: args{
			data: `{
				"Test1" : "str Val",
				"Test2" : 1
			 }`},
			want: `{"Test1":"strVal","Test2":1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canonicalize(tt.args.data)
			require.Equal(t, tt.want, got)
		})
	}
}
