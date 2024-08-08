package types

type MonthBalance struct {
	Count     int
	AvgCredit float64
	AvgDebit  float64
	Name      string
}

type TransactionsByMonth map[int]MonthBalance

type Report struct {
	TotalBalance        float64
	TransactionsByMonth TransactionsByMonth
}
