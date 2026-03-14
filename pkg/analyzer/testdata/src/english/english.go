package english

import "log/slog"

func examples() {
	slog.Info("запуск сервера")                    // want `log message must be in English`
	slog.Error("ошибка подключения к базе данных") // want `log message must be in English`
	slog.Info("starting server")                   // OK
	slog.Error("failed to connect to database")    // OK
}
