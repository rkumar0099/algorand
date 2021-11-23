package logs

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type StandardLogger struct {
	*logrus.Logger
}

var standardLogger = &StandardLogger{logrus.New()}

func init() {
	var logFilePath string = "../logs/stdout.txt"
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    5,
		MaxBackups: 10,
		MaxAge:     30,
	}

	out := os.Stdout
	mw := io.MultiWriter(out, lumberJackLogger)
	standardLogger.Formatter = &logrus.JSONFormatter{}
	standardLogger.SetOutput(mw)
}

func NewLogger() *StandardLogger {
	return standardLogger
}
