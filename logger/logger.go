package logger

import (
	"fmt"
	"time"
)

type Logger interface {
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
}

type LogLevel int

const (
	LevelInfo LogLevel = iota
	LevelWarn
	LevelError
)

type Log struct {
	Type      string        `json:"type"`
	Actor     string        `json:"actor"`
	TypeColor string        `json:"-"`
	Time      string        `json:"time"`
	Timestamp int64         `json:"timestamp"`
	Msg       string        `json:"msg"`
	MsgColor  string        `json:"-"`
	Args      []interface{} `json:"args"`
	Str       string        `json:"-"`
}

type SLogger struct {
	writers   []Writer
	color     string
	actorName string
	level     LogLevel
}

func NewLogger(actorName string, color string, level LogLevel, writers ...Writer) *SLogger {
	return &SLogger{
		writers:   writers,
		color:     color,
		actorName: actorName,
		level:     level,
	}
}

func (s *SLogger) Info(msg string, args ...interface{}) {
	log := &Log{
		Type:      "INFO",
		Actor:     s.actorName,
		TypeColor: ColorBlue,
		Time:      GetTime(),
		Timestamp: time.Now().UnixMilli(),
		Msg:       msg,
		MsgColor:  s.color,
		Args:      args,
	}
	log.Str = s.toString(log)

	if s.level <= LevelInfo {
		s.write(log)
	}
}

func (s *SLogger) Warn(msg string, args ...interface{}) {
	log := &Log{
		Type:      "WARN",
		Actor:     s.actorName,
		TypeColor: ColorOrange,
		Time:      GetTime(),
		Timestamp: time.Now().UnixMilli(),
		Msg:       msg,
		MsgColor:  s.color,
		Args:      args,
	}
	log.Str = s.toString(log)

	if s.level <= LevelWarn {
		s.write(log)
	}
}

func (s *SLogger) Error(msg string, args ...interface{}) {
	log := &Log{
		Type:      "EROR",
		Actor:     s.actorName,
		TypeColor: ColorRed,
		Time:      GetTime(),
		Timestamp: time.Now().UnixMilli(),
		Msg:       msg,
		MsgColor:  s.color,
		Args:      args,
	}
	log.Str = s.toString(log)

	if s.level <= LevelError {
		s.write(log)
	}
}

func (s *SLogger) write(l *Log) {
	for _, w := range s.writers {
		err := w.Write(l)
		if err != nil {
			fmt.Println("Error writing log", err.Error())
		}
	}
}

func argsToString(args []interface{}) []string {
	var strArgs []string
	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprintf("%v", args[i])

		var value string
		if i+1 < len(args) {
			value = fmt.Sprintf("%v", args[i+1])
		} else {
			value = "?"
		}

		strArgs = append(strArgs, key+"="+value)
	}
	return strArgs
}

func (s *SLogger) toString(log *Log) string {
	var logStr string

	logStr += FontBold + log.TypeColor + "[" + log.Type + "]" + FontNormal
	logStr += ColorWhite + " [" + log.Time + "]"
	logStr += log.MsgColor + " " + s.actorName + " "
	logStr += pad(log.Msg) + ColorWhite

	if len(log.Args) > 0 {
		logStr += fmt.Sprintf(" %v", argsToString(log.Args))
	}

	logStr += "\n"

	return logStr
}
