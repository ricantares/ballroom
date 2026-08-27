package logger

import (
	"io"
	"log"
)

type LogLevel int

const (
	Error LogLevel = iota + 1
	Warning
	Info
	Debug
)

func (t LogLevel) String() string {
	return [...]string{"Error", "Warning", "Info", "Debug"}[t-1]
}

var level LogLevel
var debugLog *log.Logger
var infoLog *log.Logger
var warningLog *log.Logger
var errorLog *log.Logger

// Inizializzazione logger. ('out' file di output del log, 'l' stringa contenente il livello di log).
func NewLog(out io.Writer, l string) {
	level = setLogLevel(l)
	debugLog = log.New(out, "DEBUG  : ", log.Ldate|log.Ltime)
	infoLog = log.New(out, "INFO   : ", log.Ldate|log.Ltime)
	warningLog = log.New(out, "WARNING: ", log.Ldate|log.Ltime)
	errorLog = log.New(out, "ERROR  : ", log.Ldate|log.Ltime)
}

func setLogLevel(l string) LogLevel {
	switch l {
	case "Debug":
		return Debug
	case "Info":
		return Info
	case "Warning":
		return Warning
	default:
		return Error
	}
}

func LogDebug(msg string) {
	if level > 3 {
		debugLog.Println(msg)
	}
}

func LogInfo(msg string) {
	if level > 2 {
		infoLog.Println(msg)
	}
}

func LogWarn(msg string) {
	if level > 1 {
		warningLog.Println(msg)
	}
}

func LogError(msg string) {
	errorLog.Println(msg)
}
