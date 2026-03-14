package zap

import "go.uber.org/zap"

func examples() {
	zap.L().Info("Starting server")                 // want `log message must start with lowercase letter`
	zap.L().Info("starting server")                 // OK
	zap.L().Error("ошибка подключения")             // want `log message must be in English`
	zap.L().Warn("server started!!!")               // want `log message must not contain special characters`
	zap.L().Info("user password: 123")              // want `log message must not contain sensitive data`
	zap.L().Info("user authenticated successfully") // OK

	zap.Info("Starting server directly")  // want `log message must start with lowercase letter`
	zap.Info("starting server directly")  // OK
	zap.Error("ошибка прямого вызова")    // want `log message must be in English`
	zap.Warn("direct call!!!")            // want `log message must not contain special characters`
	zap.Info("user password: direct")     // want `log message must not contain sensitive data`
	zap.Info("direct call authenticated") // OK
}
