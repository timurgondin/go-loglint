package sensitive

import "log/slog"

func examples() {
	slog.Info("user password: 123456")           // want `log message must not contain sensitive data`
	slog.Debug("api_key=abc123")                 // want `log message must not contain sensitive data`
	slog.Info("user authenticated successfully") // OK
	slog.Debug("api request completed")          // OK
}
