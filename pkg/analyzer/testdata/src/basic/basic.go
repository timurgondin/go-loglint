package basic

import "log/slog"

func examples() {
	slog.Info("Starting server on port 8080") // want `log message must start with lowercase letter`
	slog.Info("starting server on port 8080")
	slog.Error("ошибка подключения к базе данных") // want `log message must be in English`
	slog.Warn("server started!!!")                 // want `log message must not contain special characters`
	slog.Info("user password: 123456")             // want `log message must not contain sensitive data`
	slog.Info("user authenticated successfully")
	slog.Info("user credit_card: 1234") // want `log message must not contain sensitive data`
}
