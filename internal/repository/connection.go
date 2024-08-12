package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	secrets_service "github.com/rodrigoenzohernandez/transactions-processor/internal/services/secrets"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"

	_ "github.com/lib/pq"
)

var log = logger.GetLogger("repository")

type Secret struct {
	Password             string `json:"password"`
	DBName               string `json:"dbname"`
	Engine               string `json:"engine"`
	Port                 int    `json:"port"`
	DBInstanceIdentifier string `json:"dbInstanceIdentifier"`
	Host                 string `json:"host"`
	Username             string `json:"username"`
}

func Connect() (*sql.DB, string) {
	secret, _ := secrets_service.GetSecret("transactionsProcessorDB", "us-east-2")

	var secretData Secret
	err := json.Unmarshal([]byte(secret), &secretData)
	if err != nil {
		log.Error(fmt.Sprintf("Error parsing secret from AWS Secrets Manager: %v", err))
	}

	host := secretData.Host
	user := secretData.Username
	dbName := secretData.DBName
	password := secretData.Password
	connectTimeout := "5"
	SSLMode := "require"

	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s connect_timeout=%s sslmode=%s",
		host, user, dbName, password, connectTimeout, SSLMode)

	migrationString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=%s",
		user,
		password,
		host,
		dbName,
		SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Error connecting to the database.")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Error("Error reaching to the database.")
		panic(err)
	}

	log.Info("Successfully connected to the database.")

	return db, migrationString
}

func Disconnect(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Error("Error disconnecting from the database.")
		panic(err)
	}
}
