package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/rodrigoenzohernandez/transactions-processor/internal/models"
)

type Transaction models.Transaction

type TransactionRepoInterface interface {
	InsertMany(transactions [][]string) error
}

type TransactionRepo struct {
	DB *sql.DB
}

func NewTransactionRepo(db *sql.DB) TransactionRepoInterface {
	return &TransactionRepo{DB: db}
}

func (repo *TransactionRepo) InsertMany(transactions [][]string) error {
	// To reduce the scope it's assumed that it's just one account, so the account_id is hardcoded.
	account_id := "697ac68b-3c03-4c65-a8e1-d35e7452ba27"
	query := `INSERT INTO "dev".transactions (account_id, date, amount) VALUES `
	values := []interface{}{}
	placeholders := []string{}

	for i, txn := range transactions {
		if len(txn) != 3 {
			log.Error("invalid transaction")
		}

		date := txn[1]
		amount := txn[2]

		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			log.Error(fmt.Sprintf("invalid amount format: %v", err))

		}

		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		values = append(values, account_id, date, amountFloat)
	}

	query += strings.Join(placeholders, ", ")

	_, err := repo.DB.Exec(query, values...)
	if err != nil {
		log.Error(fmt.Sprintf("Error inserting the transactions into the database: %v", err))

		return err
	}

	log.Info("The transactions were successfully inserted into the database")

	return nil
}
