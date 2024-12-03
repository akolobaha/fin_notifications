package entity

type Target struct {
	Id                 int64
	Ticker             string
	ValuationRatio     string
	Value              float64
	FinancialReport    string
	Achieved           bool
	NotificationMethod string
}
