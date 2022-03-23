package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const (
	LogModeStdout  = 1
	LogModeLogfile = 2
	LogModeBoth    = 3
)

var (
	Logger   *logger
	levelMap = map[logLevel]string{
		DEBUG:   "debug",
		INFO:    "info",
		WARNING: "warning",
		ERROR:   "error",
		FATAL:   "fatal",
	}
)

type logger struct {
	level logLevel
	lock  *sync.Mutex
	file  *os.File
	mode  int
}

func InitLogger(mode, level int, logFile string) {
	var err error
	var file *os.File
	if mode != LogModeStdout {
		file, err = os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("open log file:%s err: %v", logFile, err)
		}
	}

	Logger = &logger{
		level: logLevel(level),
		lock:  new(sync.Mutex),
		file:  file,
		mode:  mode,
	}
}

func (l *logger) Debug(v ...interface{}) {
	l.print(DEBUG, v...)
}

func (l *logger) DebugF(format string, v ...interface{}) {
	l.printf(DEBUG, format, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.print(INFO, v...)
}

func (l *logger) InfoF(format string, v ...interface{}) {
	l.printf(INFO, format, v...)
}

func (l *logger) Warning(v ...interface{}) {
	l.print(WARNING, v...)
}

func (l *logger) WarningF(format string, v ...interface{}) {
	l.printf(WARNING, format, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.print(ERROR, v...)
}

func (l *logger) ErrorF(format string, v ...interface{}) {
	l.printf(ERROR, format, v...)
}

func (l *logger) Fatal(v ...interface{}) {
	if FATAL >= l.level {
		l.write(FATAL, fmt.Sprint(v...))
		if l.mode != LogModeLogfile {
			log.Fatal(v...)
		}
	}
}

func (l *logger) FatalF(format string, v ...interface{}) {
	if FATAL >= l.level {
		l.write(FATAL, fmt.Sprintf(format, v...))
		if l.mode != LogModeLogfile {
			log.Fatalf(format, v...)
		}
	}
}

func (l *logger) print(level logLevel, v ...interface{}) {
	if level >= l.level {
		l.write(level, fmt.Sprint(v...))
		if l.mode != LogModeLogfile {
			log.Print(v...)
		}
	}
}

func (l *logger) printf(level logLevel, format string, v ...interface{}) {
	if level >= l.level {
		l.write(level, fmt.Sprintf(format, v...))
		if l.mode != LogModeLogfile {
			log.Printf(format, v...)
		}
	}
}

func (l *logger) write(level logLevel, str string) {
	if l.file == nil || l.mode == LogModeStdout {
		return
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	// 结尾自动空格
	if len(str) == 0 || str[len(str)-1] != '\n' {
		str += "\n"
	}

	now := time.Now().Format("2006/01/02 15:04:05")
	levelStr, _ := levelMap[level]
	str = fmt.Sprintf("%s %s %s", now, levelStr, str)

	_, err := l.file.WriteString(str)
	if err != nil {
		log.Fatalf("write log file: %s, err: %v", l.file.Name(), err)
	}
}
