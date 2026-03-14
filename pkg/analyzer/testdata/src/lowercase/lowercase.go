package lowercase

import "log/slog"

func examples() {
	slog.Info("Starting server on port 8080") // want `log message must start with lowercase letter`
	slog.Error("Failed to connect")           // want `log message must start with lowercase letter`
	slog.Info("starting server on port 8080") // OK
	slog.Error("failed to connect")           // OK
}
