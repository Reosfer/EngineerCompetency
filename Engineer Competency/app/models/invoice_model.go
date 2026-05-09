package models

import "time"

type Invoice struct {
	ID          int        `json:"id"`
	BuyerId     int64      `json:"buyer_id"`
	MerchantID  int64      `json:"merchant_id"`
	AdminID     int64      `json:"admin_id"`
	InvoiceName string     `json:"invoice_name"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"`
	PaymentType string     `json:"payment_type"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type InvoiceResponseChan struct {
	Invoice *Invoice
	Error   error
}

type InvoiceListResponseChan struct {
	Invoice []Invoice
	Error   error
}

type InvoiceRequest struct {
	ID          int        `json:"invoice_id"`
	RoleID      string     `json:"role_id"`
	BuyerId     int64      `json:"buyer_id"`
	MerchantID  int64      `json:"merchant_id"`
	AdminID     int64      `json:"admin_id"`
	InvoiceName string     `json:"invoice_name"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"`
	PaymentType string     `json:"payment_type"`
	VAWallet    float64    `json:"va_wallet"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
