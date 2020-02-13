package bca

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generateSignature(t *testing.T) {
	// Test 1 & 2 based on Authentication > Signature section in https://developer.bca.co.id/documentation/#signature
	// Test 3 & 4 based on testGenerateSign() & testGenerateSign2() in  https://github.com/odenktools/php-bca/blob/develop/test/unit/bcaConstructorTest.php

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
		{name: "official bca/POST Example", args: args{
			apiSecret:   "22a2d25e-765d-41e1-8d29-da68dcb5698b",
			method:      http.MethodPost,
			path:        "/banking/corporates/transfers",
			accessToken: "lIWOt2p29grUo59bedBUrBY3pnzqQX544LzYPohcGHOuwn8AUEdUKS",
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
		{name: "official bca/GET Example", args: args{
			apiSecret:   "22a2d25e-765d-41e1-8d29-da68dcb5698b",
			method:      http.MethodGet,
			path:        "/banking/v2/corporates/BCAAPI2016/accounts/0201245680/statements?StartDate=2016-09-01&EndDate=2016-09-01",
			accessToken: "lIWOt2p29grUo59bedBUrBY3pnzqQX544LzYPohcGHOuwn8AUEdUKS",
			requestBody: "",
			timestamp:   "2016-02-03T10:00:00.000+07:00",
		},
			wantSign:      "3ac124303746d222387d4398dddf33201a384aa22137aa08f4d9843c6f467a48",
			wantStrToSign: "GET:/banking/v2/corporates/BCAAPI2016/accounts/0201245680/statements?EndDate=2016-09-01&StartDate=2016-09-01:lIWOt2p29grUo59bedBUrBY3pnzqQX544LzYPohcGHOuwn8AUEdUKS:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855:2016-02-03T10:00:00.000+07:00",
			wantErr:       false,
		},
		{name: "php-bca/testGenerateSign/GET Example", args: args{
			apiSecret:   "9db65b91-01ff-46ec-9274-3f234b677450",
			method:      http.MethodGet,
			path:        "/banking/v2/corporates/corpid/accounts/0063001004",
			accessToken: "NopUsBuSbT3eNrQTfcEZN2aAL52JT1SlRgoL1MIslsX5gGIgv4YUf",
			requestBody: "",
			timestamp:   "2017-09-30T22:03:35.800+07:00",
		},
			wantSign:      "761eaec0e544c9cf5010b406ade39228ab182401e57f17fc54b9daa5ad99d0d6",
			wantStrToSign: "GET:/banking/v2/corporates/corpid/accounts/0063001004:NopUsBuSbT3eNrQTfcEZN2aAL52JT1SlRgoL1MIslsX5gGIgv4YUf:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855:2017-09-30T22:03:35.800+07:00",
			wantErr:       false,
		},
		// TODO: STILL FAILED
		// {name: "php-bca/testGenerateSign2/GET Example", args: args{
		// 	apiSecret:   "9db65b91-01ff-46ec-9274-3f234b677450",
		// 	method:      http.MethodGet,
		// 	path:        "/banking/v2/corporates/corpid/accounts/0063001004",
		// 	accessToken: "NopUsBuSbT3eNrQTfcEZN2aAL52JT1SlRgoL1MIslsX5gGIgv4YUf",
		// 	requestBody: `
		// 	{
		// 		"Amount" : "100000.00",
		// 		"BeneficiaryAccountNumber" : "8329389",
		// 		"CorporateID" : "8293489283499",
		// 		"CurrencyCode" : "idr",
		// 		"ReferenceID" : "",
		// 		"Remark1" : "Ini adalah remark1",
		// 		"Remark2" : "Ini adalah remark2",
		// 		"SourceAccountNumber" : "09202990",
		// 		"TransactionDate" : "2019-02-30T22:03:35.800+07:00",
		// 		"TransactionID" : "0020292"
		// 	}
		// 	`,
		// 	timestamp: "2017-09-30T22:03:35.800+07:00",
		// },
		// 	wantSign:      "1878f0eedcd93ff53054c8fc9ea271a29c99ea2f752f636c1cc765948009a90b",
		// 	wantStrToSign: "GET:/banking/v2/corporates/corpid/accounts/0063001004:NopUsBuSbT3eNrQTfcEZN2aAL52JT1SlRgoL1MIslsX5gGIgv4YUf:4a24a20ceb436d69bd344902a71e9bdf3a45d11efac1754b48015fb2a291b3df:2017-09-30T22:03:35.800+07:00",
		// 	wantErr:       false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSign, gotStrToSign, err := GenerateSignature(tt.args.apiSecret, tt.args.method, tt.args.path, tt.args.accessToken, tt.args.requestBody, tt.args.timestamp)
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
