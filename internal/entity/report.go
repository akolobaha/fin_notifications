package entity

type Report struct {
	TargetId            int64
	Ticker              string
	ValuationRatio      string
	Value               float64
	FinancialReport     string
	NotificationMethod  string
	NotificationSubject string
	NotificationText    string
}

func NewReport(targetUser TargetUser, subject string, text string) *Report {
	return &Report{
		TargetId:            targetUser.Target.Id,
		Ticker:              targetUser.Target.Ticker,
		ValuationRatio:      targetUser.Target.ValuationRatio,
		Value:               targetUser.Target.Value,
		FinancialReport:     targetUser.Target.FinancialReport,
		NotificationMethod:  targetUser.Target.NotificationMethod,
		NotificationSubject: subject,
		NotificationText:    text,
	}
}
