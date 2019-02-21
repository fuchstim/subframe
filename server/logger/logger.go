package logger

import (
	"errors"
	"fmt"
	"github.com/bclicn/color"
	"os"
	"strconv"
	"time"
)

//LogtypeInfo defines Info-Level Logs
const LogtypeInfo = 1

//LogtypeWarn defines Warn-Level Logs
const LogtypeWarn = 2

//LogtypeError defines Error-Level Logs
const LogtypeError = 3

//LogtypeFatal defines Fatal-Level Logs
const LogtypeFatal = 4

//Initialized shows whether the Logger has been initialized
var Initialized = false

var logtypeDescriptions = [...]string{"INFO", "WARN", "ERROR", "FATAL"}
var logFile *os.File
var logLogger = Logger{"logger/Logger"}

//LogPath is the Path at which log files reside
var LogPath string

//ColorizedLogs turns on or off colorized realtime log output
var ColorizedLogs bool

var logQueue []string

//Init initializes the Logger
func Init() {
	logLogger.Info("Initializing Logger...")
	logFile = makeLogFile()
	Initialized = true
	logLogger.Info("Initialized Logger")
}

//Close closes the Logger
func Close() {
	logLogger.Info("Closing Logger...")
	logFile.Close()
	Initialized = false
	logLogger.Info("Closed Logger.")
}

//Logger creates a new Logger for a specific context
type Logger struct {
	Prefix string
}

//Log holds relevant information about a log element
type Log struct {
	Time    time.Time
	Prefix  string
	Type    int
	Message string
}

//Info logs an info-message
func (l Logger) Info(message string) {
	mainLogger(Log{time.Now(), l.Prefix, LogtypeInfo, message})
}

//Warn logs a warn-message
func (l Logger) Warn(message string) {
	mainLogger(Log{time.Now(), l.Prefix, LogtypeWarn, message})
}

//Error logs an error-message
func (l Logger) Error(message string) {
	mainLogger(Log{time.Now(), l.Prefix, LogtypeError, message})
}

//Fatal logs a fatal message and panics
func (l Logger) Fatal(message string) {
	mainLogger(Log{time.Now(), l.Prefix, LogtypeFatal, message})
	panic(errors.New(message))
}

var lastPrefix string

func mainLogger(l Log) {
	if lastPrefix != "" && lastPrefix != l.Prefix {
		//Log empty line when changing Prefixes / Contexts
		logToCLI("")
		logToFile("")
	}

	if !Initialized || !ColorizedLogs {
		logToCLI(formatLogLine(l))
	} else {
		logToCLI(formatLogLineCLI(l))
	}
	logToFile(formatLogLine(l))
	lastPrefix = l.Prefix
}

func logToCLI(logLine string) {
	fmt.Println(logLine)
}

func logToFile(logLine string) {
	if Initialized {
		//Write buffer to file
		if len(logQueue) > 0 {

			l := Log{
				time.Now(),
				"logger/Logger",
				LogtypeInfo,
				"Writing logQueue (" + strconv.Itoa(len(logQueue)) + " Elements) to logfile...",
			}
			if ColorizedLogs {
				logToCLI(formatLogLineCLI(l))
			} else {
				logToCLI(formatLogLine(l))
			}

			for _, value := range logQueue {
				logFile.WriteString(value + "\n")
			}
			logFile.Sync()
			logQueue = logQueue[:0]

			l = Log{
				time.Now(),
				"logger/Logger",
				LogtypeInfo,
				"Wrote logQueue to logfile...",
			}

			if ColorizedLogs {
				logToCLI(formatLogLineCLI(l))
			} else {
				logToCLI(formatLogLine(l))
			}
		}
		logFile.WriteString(logLine + "\n")
		logFile.Sync()
		return
	}
	//If logger is not yet initialized, cache log lines to buffer
	logQueue = append(logQueue, logLine)
}

func makeLogFile() (logFile *os.File) {
	file, err := os.Create(LogPath + "/" + time.Now().Format("Mon-Jan-2-2006-15-04-05.log"))
	if err != nil {
		panic(err)
	}
	return file
}

func formatTime(t time.Time) (formatted string) {
	return "[" + t.Format("Mon Jan 2 2006 15:04:05") + "]"
}

func formatPrefix(prefix string) (formatted string) {
	return "(" + prefix + ")"
}

func formatLogType(logType int) (formatted string) {
	return "[" + logtypeDescriptions[logType-1] + "]"
}

func formatLogLine(l Log) (line string) {
	//Format like:
	//[Mon Jan 1 12:13:14 2019] [INFO] [logger/Init] Initialized Logger.
	return formatTime(l.Time) + " " + formatLogType(l.Type) + " " + formatPrefix(l.Prefix) + " " + l.Message
}

func formatTimeCLI(t time.Time) (formatted string) {
	return "[" + color.Green(t.Format("Mon Jan 2 2006 15:04:05")) + "]"
}

func formatPrefixCLI(prefix string) (formatted string) {
	return "(" + color.BBlue(prefix) + ")"
}

func formatLogTypeCLI(logType int) (formatted string) {
	switch logType {
	case LogtypeInfo:
		return "[" + color.BGreen(logtypeDescriptions[logType-1]) + "]"
	case LogtypeWarn:
		return "[" + color.BYellow(logtypeDescriptions[logType-1]) + "]"
	case LogtypeError:
		return "[" + color.BLightRed(logtypeDescriptions[logType-1]) + "]"
	case LogtypeFatal:
		return "[" + color.BRed(logtypeDescriptions[logType-1]) + "]"
	}
	return "[" + logtypeDescriptions[logType-1] + "]"
}

func formatLogLineCLI(l Log) (line string) {
	//Format like:
	//[Mon Jan 1 12:13:14 2019] [INFO] [logger/Init] Initialized Logger.
	return formatTimeCLI(l.Time) + " " + formatLogTypeCLI(l.Type) + " " + formatPrefixCLI(l.Prefix) + " " + l.Message
}
