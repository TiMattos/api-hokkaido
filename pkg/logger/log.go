package logger

import (
	"time"

	"go.uber.org/zap"
)

func GravarLog(mensagem string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("Erro no processamento",
		// Structured context as loosely typed key-value pairs.
		"message", mensagem,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Gravando log de processamento: %s", mensagem)
}
