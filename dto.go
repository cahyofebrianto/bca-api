package bca

import validation "github.com/go-ozzo/ozzo-validation"

// === AUTH ===

// AuthToken represents response of BCA OAuth 2.0 response message
type AuthToken struct {
	Error
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// === ERROR ===

// Error represent BCA error response messsage
type Error struct {
	ErrorCode    string
	ErrorMessage ErrorLang
}

// ErrorLang represent BCA error response message language
type ErrorLang struct {
	Indonesian string
	English    string
}

// === BCA ===

// AccountBalance represents account balance information
type AccountBalance struct {
	AccountNumber    string
	Currency         string  `json:",omitempty"`
	Balance          float64 `json:",string"`
	AvailableBalance float64 `json:",string"`
	FloatAmount      float64 `json:",string"`
	HoldAmount       float64 `json:",string"`
	Plafon           float64 `json:",string"`
	Indonesian       string  `json:",omitempty"`
	English          string  `json:",omitempty"`
}

// BalanceInfoRequest represents account balance information request message
type BalanceInfoRequest struct {
	AccountNumber string
}

// BalanceInfoResponse represents account balance information response message
type BalanceInfoResponse struct {
	Error
	AccountDetailDataSuccess []AccountBalance `json:",omitempty"`
	AccountDetailDataFailed  []AccountBalance `json:",omitempty"`
}

// AccountStatement represents account statement information
type AccountStatement struct {
	TransactionDate   string
	BranchCode        string
	TransactionType   string
	TransactionAmount float64 `json:",string"`
	TransactionName   string
	Trailer           string
}

// AccountStatementResponse represents account statement response message
type AccountStatementResponse struct {
	Error
	StartDate    string
	EndDate      string
	Currency     string
	StartBalance float64 `json:",string"`
	Data         []AccountStatement
}

// FundTransferRequest represents fund transfer request message
type FundTransferRequest struct {
	CorporateID              string
	SourceAccountNumber      string
	TransactionID            string
	TransactionDate          string
	ReferenceID              string
	CurrencyCode             string
	Amount                   float64 `json:",string"`
	BeneficiaryAccountNumber string
	Remark1                  string
	Remark2                  string
}

// FundTransferResponse represents fund transfer response message
type FundTransferResponse struct {
	Error
	TransactionID   string
	TransactionDate string
	ReferenceID     string
	Status          string
}

// FundTransferDomesticRequest represents fund transfer request message
type FundTransferDomesticRequest struct {
	TransactionID            string
	TransactionDate          string
	ReferenceID              string
	SourceAccountNumber      string
	BeneficiaryAccountNumber string
	BeneficiaryBankCode      string
	BeneficiaryName          string
	Amount                   float64 `json:",string"`
	TransferType             string
	BeneficiaryCustType      string
	BeneficiaryCustResidence string
	CurrencyCode             string
	Remark1                  string
	Remark2                  string
}

// FundTransferDomesticResponse represents fund transfer response message
type FundTransferDomesticResponse struct {
	Error
	TransactionID   string
	TransactionDate string
	ReferenceID     string
	PPUNumber       string
	Status          string
}

// InquiryBillRequest represents VA inquiry bill message
type InquiryBillRequest struct {
	CompanyCode     string
	CustomerNumber  string
	RequestID       string
	ChannelType     string
	TransactionDate string
	AdditionalData  string
}

func (m InquiryBillRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CustomerNumber, validation.Required),
		validation.Field(&m.CompanyCode, validation.Required),
		validation.Field(&m.ChannelType, validation.Required),
		validation.Field(&m.RequestID, validation.Required),
		validation.Field(&m.TransactionDate, validation.Required),
	)
}
func (m InquiryBillRequest) ValidateTransDate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.TransactionDate, validation.Date("02/01/2006 15:04:05")),
	)
}

// ReasonMessage ...
type ReasonMessage struct {
	Indonesian string
	English    string
}

// DetailBill ...
type DetailBill struct {
	BillDescription ReasonMessage
	BillAmount      string
	BillNumber      string
	BillSubCompany  string
}

type DetailBillPayment struct {
	BillNumber 		string
	Status			string
	Reason 			ReasonMessage
}

// InquiryBillSingleResponse ...
type InquiryBillSingleResponse struct {
	CompanyCode		string
	CustomerNumber 	string
	RequestID      	string
	InquiryStatus  	string
	InquiryReason  	ReasonMessage
	CustomerName   	string
	CurrencyCode   	string
	TotalAmount    	string
	SubCompany     	string
	DetailBills    	[]DetailBill
	FreeTexts 		[]ReasonMessage
	AdditionalData	string
}

// PaymentBillRequest ...
type PaymentBillRequest struct {
	CompanyCode     string
	CustomerNumber  string
	RequestID       string
	ChannelType     string
	CustomerName    string
	CurrencyCode    string
	PaidAmount      string
	TotalAmount     string
	SubCompany      string
	TransactionDate string
	Reference       string
	DetailBills     []DetailBill
	FlagAdvice      string
	AdditionalData  string
}

func (m PaymentBillRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CompanyCode, validation.Required),
		validation.Field(&m.CustomerNumber, validation.Required),
		validation.Field(&m.RequestID, validation.Required),
		validation.Field(&m.ChannelType, validation.Required),
		validation.Field(&m.TransactionDate, validation.Required),
		validation.Field(&m.PaidAmount, validation.Required),
	)
}
func (m PaymentBillRequest) ValidateTransDate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.TransactionDate, validation.Date("02/01/2006 15:04:05")),
	)
}

// PaymentBillResponse ...
type PaymentBillResponse struct {
	CompanyCode       string
	CustomerNumber    string
	RequestID         string
	PaymentFlagStatus string
	PaymentFlagReason ReasonMessage
	CurrencyCode      string
	PaidAmount        string
	TotalAmount       string
	TransactionDate   string
	DetailBills       []DetailBillPayment
	FreeTexts         []ReasonMessage
	AdditionalData    string
}

// Authentication mostly is used as embedded struct for Fire API request
type Authentication struct {
	CorporateID string
	AccessCode  string
	BranchCode  string
	UserID      string
	LocalID     string
}

// InquiryAccountRequestBeneficiaryDetails is beneficiary details of inquiry account request
type InquiryAccountRequestBeneficiaryDetails struct {
	BankCodeType  string
	BankCodeValue string
	AccountNumber string
}

// InquiryAccountResponseBeneficiaryDetails is beneficiary details of inquiry account response
type InquiryAccountResponseBeneficiaryDetails struct {
	ServerBeneAccountName string
}

// InquiryAccountRequest represents inquiry account request message
type InquiryAccountRequest struct {
	Authentication     Authentication
	BeneficiaryDetails InquiryAccountRequestBeneficiaryDetails
}

// InquiryAccountResponse represents inquiry account response message
type InquiryAccountResponse struct {
	Error
	BeneficiaryDetails InquiryAccountResponseBeneficiaryDetails
	StatusTransaction  string
	StatusMessage      string
}
