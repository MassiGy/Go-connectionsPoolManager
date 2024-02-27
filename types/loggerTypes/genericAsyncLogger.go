package loggertypes

import (
	"fmt"
	"time"
)

type GenericAsyncLogger struct {

	// config object, we can also embbed the type instead
	Config AsyncLoggerConfig
}

func (genericAsyncLogger GenericAsyncLogger) GetConfig() AsyncLoggerConfig {
	return genericAsyncLogger.Config
}

// this will act as a middelware between the users of the channel and the logger
// internals. Hopefully this will make the users only use the channel as a write-to channel
func (genericAsyncLogger *GenericAsyncLogger) GetAsyncLoggerHandle() chan<- []byte {
	return genericAsyncLogger.Config.buffer
}

func (genericAsyncLogger *GenericAsyncLogger) SetAutoFlush() error {
	genericAsyncLogger.Config.isAutoFlush = true
	return nil
}

// this won't check if the logger buffer is closed,
func (genericAsyncLogger GenericAsyncLogger) EnQueue(msg []byte) {
	genericAsyncLogger.Config.buffer <- msg
}

func (genericAsyncLogger *GenericAsyncLogger) Listen() error {
	var (
		err error
	)

	// for every tick, flush to logger sink
	for tick := range genericAsyncLogger.Config.tick.C {

		// pass the time instant to flush
		err = genericAsyncLogger.Flush(tick)
		if err != nil {
			return err
		}
	}

	return nil
}

func (genericAsyncLogger *GenericAsyncLogger) Flush(timeStamp time.Time) error {

	fmt.Fprintf(genericAsyncLogger.Config.Sink, "Start of tick ===============\n")

	// while true
	for {

		// either consume the msg, or quit
		select {
		case msg := <-genericAsyncLogger.Config.buffer:
			{
				fmt.Fprintf(genericAsyncLogger.Config.Sink, "\n[%d:%d:%d]\t", timeStamp.Hour(), timeStamp.Minute(), timeStamp.Second())
				fmt.Fprintf(genericAsyncLogger.Config.Sink, "@%s:\t", genericAsyncLogger.Config.name)
				fmt.Fprintf(genericAsyncLogger.Config.Sink, "(%s)\t", severityLevels[genericAsyncLogger.Config.severityLevel])
				fmt.Fprintf(genericAsyncLogger.Config.Sink, "%s", msg)
			}
		case <-time.After(genericAsyncLogger.Config.flushTimeOut):
			fmt.Fprintf(genericAsyncLogger.Config.Sink, "End of tick ===============\n")
			return nil
		}
	}
}

func (genericAsyncLogger *GenericAsyncLogger) Close() error {

	// we need to close the buffer
	close(genericAsyncLogger.Config.buffer)

	return nil
}
