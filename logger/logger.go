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
	INFO LogLevel = iota
	WARN
	ERROR
)

type Log struct {
	Type      string        `json:"type"`
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
		TypeColor: BLUE,
		Time:      GetTime(),
		Timestamp: time.Now().UnixMilli(),
		Msg:       msg,
		MsgColor:  s.color,
		Args:      args,
	}
	log.Str = s.toString(log)

	if s.level <= INFO {
		s.write(log)
	}
}

func (s *SLogger) Warn(msg string, args ...interface{}) {
	log := &Log{
		Type:      "WARN",
		TypeColor: ORANGE,
		Time:      GetTime(),
		Timestamp: time.Now().UnixMilli(),
		Msg:       msg,
		MsgColor:  s.color,
		Args:      args,
	}
	log.Str = s.toString(log)

	if s.level <= WARN {
		s.write(log)
	}
}

func (s *SLogger) Error(msg string, args ...interface{}) {
	log := &Log{
		Type:      "ERROR",
		TypeColor: RED,
		Time:      GetTime(),
		Timestamp: time.Now().UnixMilli(),
		Msg:       msg,
		MsgColor:  s.color,
		Args:      args,
	}
	log.Str = s.toString(log)

	if s.level <= ERROR {
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

	logStr += BOLD + log.TypeColor + "[" + log.Type + "]" + NORMAL
	logStr += WHITE + " [" + log.Time + "]"
	logStr += log.MsgColor + " " + s.actorName + " "
	logStr += pad(log.Msg) + WHITE

	if len(log.Args) > 0 {
		logStr += fmt.Sprintf(" %v", argsToString(log.Args))
	}

	logStr += "\n"

	return logStr
}
