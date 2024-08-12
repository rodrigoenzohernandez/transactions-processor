package models

type Transaction struct {
	ID        string  `json:"id"`
	AccountID string  `json:"account_id" binding:"required"`
	Date      string  `json:"date"  binding:"required"`
	Amount    float64 `json:"amount"  binding:"required"`
}
