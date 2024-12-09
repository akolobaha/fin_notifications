package log

import (
	"fin_notifications/internal/monitoring"
	"fmt"
	"log/slog"
)

func Error(additionalMessage string, err error) {
	if err != nil {
		msg := fmt.Sprintf("%s: %s", additionalMessage, err.Error())
		monitoring.NotificationsErrorCount.WithLabelValues(msg).Inc()
		slog.Error(msg)
	}
}

func Info(message string) {
	monitoring.NotificationsSuccessCount.WithLabelValues(message).Inc()
	slog.Info(message)
}
