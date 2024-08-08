package types

type MonthBalance struct {
	Count     int
	AvgCredit float64
	AvgDebit  float64
}

type Report struct {
	TotalBalance        float64
	TransactionsByMonth map[string]MonthBalance
}
