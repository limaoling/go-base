package utils

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

// InitLogger initializes Gin-based logging and writes logs to both file and console.
func InitLogger(logPath string) error {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(f, os.Stderr)

	infoLogger = log.New(gin.DefaultWriter, "", log.LstdFlags)
	errorLogger = log.New(gin.DefaultErrorWriter, "", log.LstdFlags)
	return nil
}

func Infof(format string, args ...any) {
	infoLogger.Printf(format, args...)
}

func Warnf(format string, args ...any) {
	infoLogger.Printf(format, args...)
}

func Errorf(format string, args ...any) {
	errorLogger.Printf(format, args...)
}
