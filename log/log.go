package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gitlab.trading.imc.intra/stampede/dock2box-ng/common/color"
)

var (
	logLevel   Level
	logInColor = true
)

func SetLogLevel(l Level) {
	logLevel = l
}

func SetLogLevelString(s string) error {
	for k, v := range levelToName {
		if strings.ToLower(s) == strings.ToLower(v) {
			logLevel = k
			return nil
		}
	}
	return fmt.Errorf("unknown log level: %s", s)
}

func LogLevel() Level {
	return logLevel
}

func SetMultiOutput(writers ...io.Writer) {
	log.SetOutput(io.MultiWriter(writers...))
}

func NoColor() {
	logInColor = false
}

func Print(v ...interface{}) {
	log.Print(v...)
}

func Println(v ...interface{}) {
	log.Println(v...)
}

func Printf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

func Debug(v ...interface{}) {
	Log(LevelDebug, fmt.Sprint(v...))
}

func Debugln(v ...interface{}) {
	Log(LevelDebug, fmt.Sprintln(v...))
}

func Debugf(f string, v ...interface{}) {
	Log(LevelDebug, fmt.Sprintf(f, v...))
}

func Info(v ...interface{}) {
	Log(LevelInfo, fmt.Sprint(v...))
}

func Infoln(v ...interface{}) {
	Log(LevelInfo, fmt.Sprintln(v...))
}

func Infof(f string, v ...interface{}) {
	Log(LevelInfo, fmt.Sprintf(f, v...))
}

func Warn(v ...interface{}) {
	Log(LevelWarn, fmt.Sprint(v...))
}

func Warnln(v ...interface{}) {
	Log(LevelWarn, fmt.Sprintln(v...))
}

func Warnf(f string, v ...interface{}) {
	Log(LevelWarn, fmt.Sprintf(f, v...))
}

func Error(v ...interface{}) {
	Log(LevelError, fmt.Sprint(v...))
}

func Errorln(v ...interface{}) {
	Log(LevelError, fmt.Sprintln(v...))
}

func Errorf(f string, v ...interface{}) {
	Log(LevelError, fmt.Sprintf(f, v...))
}

func Fatal(v ...interface{}) {
	Log(LevelCrit, fmt.Sprint(v...))
	os.Exit(1)
}

func Fataln(v ...interface{}) {
	Log(LevelCrit, fmt.Sprintln(v...))
	os.Exit(1)
}

func Fatalf(f string, v ...interface{}) {
	Log(LevelCrit, fmt.Sprintf(f, v...))
	os.Exit(1)
}

func Log(l Level, m string) {
	if l < logLevel {
		return
	}

	if logInColor {
		log.Printf("%s%s%s", levelToColor[l], m, color.Reset)
	} else {
		log.Print(m)
	}
}
