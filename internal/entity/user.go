package entity

type User struct {
	ID       int64
	Name     string
	Email    string
	Telegram string
}

type TargetUser struct {
	Target      Target
	User        User
	ResultValue float64
}
