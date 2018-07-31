package log

import "github.com/mickep76/grpc-exec-example/color"

type Level int

const (
	LevelNone = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCrit
)

var (
	levelToName = map[Level]string{
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARNING",
		LevelError: "ERROR",
		LevelCrit:  "CRITICAL",
	}
	levelToColor = map[Level]color.Code{
		LevelDebug: color.Cyan,
		LevelInfo:  color.LightGray,
		LevelWarn:  color.Yellow,
		LevelError: color.LightRed,
		LevelCrit:  color.Red,
	}
)

func (l Level) String() string {
	return levelToName[l]
}
