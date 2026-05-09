package models

import "time"

type Wallet struct {
	ID           int        `json:"id"`
	UserID       int64      `json:"user_id"`
	InvoiceID    int64      `json:"invoice_id"`
	CurrentSaldo float64    `json:"current_saldo"`
	RecordSaldo  float64    `json:"record_saldo"`
	Action       string     `json:"action"`
	Status       int        `json:"status"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type WalletResponseChan struct {
	Wallet *Wallet
	Error  error
}

type WalletListResponseChan struct {
	Wallet []Wallet
	Error  error
}
