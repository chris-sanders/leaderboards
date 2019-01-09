package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type logStruct struct {
	*log.Logger
	TraceLevel log.Level
	DebugLevel log.Level
	InfoLevel  log.Level
	WarnLevel  log.Level
	ErrorLevel log.Level
	FatalLevel log.Level
	PanicLevel log.Level
}

func (l *logStruct) Close() {
	logFile.Close()
}

var Log = &logStruct{
	log.New(),
	log.TraceLevel,
	log.DebugLevel,
	log.InfoLevel,
	log.WarnLevel,
	log.ErrorLevel,
	log.FatalLevel,
	log.PanicLevel,
}

var logFile *os.File

func SetupFile() {
	logFile, err := os.OpenFile("leaderboards.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Log.Fatal(err)
	}
	Log.SetOutput(logFile)
}
