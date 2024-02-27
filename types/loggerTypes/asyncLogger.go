package loggertypes

import (
	"io"
	"time"
)

// make the severityLevel type private, since this is the
// underlying type for our enum. This will prevent users from
// defining other entries of the same nature as our enums
type severityLevel int
type severityDict map[severityLevel]string

const (
	DEBUG = iota // starts at 0
	INFO
	WARNNING
	DANGER
	CRITICAL
)

var severityLevels severityDict = severityDict{
	DEBUG:    "DEBUG",
	INFO:     "INFO",
	WARNNING: "WARNNING",
	DANGER:   "DANGER",
	CRITICAL: "CRITICAL",
}

type AsyncLoggerConfig struct {

	// set the name of the logger
	name string

	// set the sink where to dump the logs so as we can pass it to fprintf
	Sink io.Writer

	// set the Severitylevel of the logger
	severityLevel severityLevel

	// set a ticker for autoflush behavior
	tick time.Ticker

	buffer       chan []byte
	isAutoFlush  bool
	flushTimeOut time.Duration
}

func (conf AsyncLoggerConfig) WithLoggerName(name string) AsyncLoggerConfig {
	conf.name = name
	return conf
}
func (conf AsyncLoggerConfig) WithLoggerSeverity(severity severityLevel) AsyncLoggerConfig {
	conf.severityLevel = severity
	return conf
}

func (conf AsyncLoggerConfig) WithTimeTick(tick time.Ticker) AsyncLoggerConfig {
	conf.tick = tick
	return conf
}

func (conf AsyncLoggerConfig) WithSink(sink io.Writer) AsyncLoggerConfig {
	conf.Sink = sink
	return conf
}

func (conf AsyncLoggerConfig) WithBuffer(buffer chan []byte) AsyncLoggerConfig {
	conf.buffer = buffer
	return conf
}
func (conf AsyncLoggerConfig) WithAutoFlushSetTo(isAutoFlush bool) AsyncLoggerConfig {
	conf.isAutoFlush = isAutoFlush
	return conf
}
func (conf AsyncLoggerConfig) WithFlushTimeOut(flushTimeOut time.Duration) AsyncLoggerConfig {
	conf.flushTimeOut = flushTimeOut
	return conf
}

type AsyncLogger interface {
	GetConfig() AsyncLoggerConfig

	GetAsyncLoggerHandle() chan<- []byte
	SetAutoFlush() error
	Listen() error
	Flush(time.Time) error
	EnQueue([]byte)
	Close() error
}

// This will be useful for anything that only
// wants to dump in msg into the logger
type LogQueuer interface {
	EnQueue([]byte)
}
