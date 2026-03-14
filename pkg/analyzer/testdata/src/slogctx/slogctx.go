package slogctx

import (
	"context"
	"log/slog"
)

func examples() {
	ctx := context.Background()

	slog.InfoContext(ctx, "Starting server")     // want `log message must start with lowercase letter`
	slog.WarnContext(ctx, "server started!!!")   // want `log message must not contain special characters`
	slog.ErrorContext(ctx, "user password: 123") // want `log message must not contain sensitive data`
	slog.DebugContext(ctx, "request processed")  // OK
	slog.InfoContext(ctx, "starting server")     // OK
}
