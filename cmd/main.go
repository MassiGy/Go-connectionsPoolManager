package main

import (
	connTypes "connectionsPoolManager/types/connectionsTypes"
	loggerTypes "connectionsPoolManager/types/loggerTypes"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {

	logger, file := setupToFileLogger("build/logs/dump.logs")
	defer file.Close()
	// listen in background
	go func() { logger.Listen() }()

	// init a connections pool
	connectionsLimitRate := 8
	connectionsPool := &connTypes.HttpConnectionPool{
		Connections: make([]connTypes.HttpConnection, 0, connectionsLimitRate),
	}
	connectionsPool.SetPoolSize(connectionsLimitRate) // pool size should'nt be exposed
	logger.EnQueue([]byte(fmt.Sprintf("Connections pool size set to %d.\n", connectionsLimitRate)))

	// init a connection pool manager
	connectionsPoolManager := &connTypes.HttpConnectionPoolManager{}
	connectionsPoolManager.ConnectionsPool = connectionsPool

	logger.EnQueue([]byte("Initialization of connection pools done.\n"))

	logger.EnQueue([]byte("About to create a connection.\n"))

	connection := setupHttpConnection(time.Minute, 1)
	connectionsPoolManager.RegisterConnection(connection)

	logger.EnQueue([]byte(fmt.Sprintf("Connection registred with id = %d.\n", connection.GetId())))

	logger.EnQueue([]byte("Connections Pool monitoring time cycle set to 3s.\n"))
	logger.EnQueue([]byte("Connections Pool monitoring time period set to 70s.\n"))

	// this works since logger satisfies LogQueuer interface
	connectionsPoolManager.SetLoggingHandler(logger)

	// make sure that period % interval != 0, (otherwise it might not end)
	connectionsPoolManager.Monitor(70*time.Second, 3*time.Second)
}

func setupToFileLogger(targetFileName string) (*loggerTypes.GenericAsyncLogger, *os.File) {
	file, err := os.Create(targetFileName)
	check(err)

	config := loggerTypes.AsyncLoggerConfig{}.
		WithLoggerName("simpleLogger").
		WithLoggerSeverity(loggerTypes.INFO).
		WithTimeTick(*time.NewTicker(time.Second)).
		WithSink(file).
		WithBuffer(make(chan []byte, 3)).
		WithAutoFlushSetTo(true).
		WithFlushTimeOut(100 * time.Millisecond)

	logger := &loggerTypes.GenericAsyncLogger{
		Config: config,
	}

	return logger, file
}

func setupHttpConnection(ttl time.Duration, id int) *connTypes.HttpConnection {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))

	conn := &connTypes.HttpConnection{}
	conn.SetId(1)
	conn.SetContext(ctx)

	return conn
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
