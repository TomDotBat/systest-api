package log

import "math"

type Level uint16

const (
	Off   Level = 0
	Fatal Level = 100
	Error Level = 200
	Warn  Level = 300
	Info  Level = 400
	Debug Level = 500
	All   Level = math.MaxUint16
)

func (level Level) String() string {
	switch level {
	case Fatal:
		return "FATAL"
	case Error:
		return "ERROR"
	case Warn:
		return "WARN"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
