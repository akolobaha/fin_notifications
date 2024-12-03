package email

import (
	"fin_notifications/internal/entity"
	"fmt"
)

func GenerateEmailText(user entity.TargetUser) (subject string, text string) {
	text = fmt.Sprintf("Цель %s по эмитенту %s достигнута:\nцель - %f, последнее значение - %f",
		getRatioText(user.Target.ValuationRatio),
		user.Target.Ticker,
		user.Target.Value,
		user.ResultValue,
	)

	subject = fmt.Sprintf("Цель по %s достигнута", user.Target.Ticker)

	return subject, text
}

func getRatioText(ratio string) string {
	switch ratio {
	case "pbv":
		return "P / Bv"
	case "pe":
		return "P / E"
	case "ps":
		return "P / S"
	case "price":
		return "Цена за акцию"
	default:
		return ""
	}

}
