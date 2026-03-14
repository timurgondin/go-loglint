package specialchars

import "log/slog"

func examples() {
	slog.Info("server started!!!")    // want `log message must not contain special characters`
	slog.Warn("server started!!!")    // want `log message must not contain special characters`
	slog.Info("server started")       // OK
	slog.Warn("something went wrong") // OK
}
