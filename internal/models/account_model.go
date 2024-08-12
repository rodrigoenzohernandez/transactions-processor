package models

type Account struct {
	ID      string `json:"id"`
	Address string `json:"address"  binding:"required"`
}
