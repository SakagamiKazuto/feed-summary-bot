package logger

import "go.uber.org/zap"

var LOG *zap.Logger

func NewLogger() {
	LOG, _ = zap.NewDevelopment()
}
