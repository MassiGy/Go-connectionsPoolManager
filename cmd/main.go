package main

import (
	"connectionsPoolManager/types"
	"context"
	"fmt"
	"time"
)

var connectionsLimitRate int
var connectionsPool types.ConnectionsPool
var connectionsPoolManager types.ConnectionsPoolManager
var loggingHandler chan []byte

func main() {
	// set the limit of connection
	connectionsLimitRate = 8

	// init a connections pool
	connectionsPool = &types.HttpConnectionPool{
		Connections: make([]types.HttpConnection, connectionsLimitRate),
	}
	connectionsPool.SetPoolSize(connectionsLimitRate)

	// init a connection pool manager
	connectionsPoolManager = &types.HttpConnectionPoolManager{}

	loggingHandler = make(chan []byte, 3) // @todo: setup an actual logger
	go func() {
		// consume the logs to prevent blocking
		fmt.Println(<-loggingHandler)
	}()
	connectionsPoolManager.SetLoggingHandler(loggingHandler)

	loggingHandler <- []byte("initialization of connection pools done.")

	loggingHandler <- []byte("about to create a connection.")

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))

	conn := &types.HttpConnection{}
	conn.SetId(1)
	conn.SetContext(ctx)

	connectionsPoolManager.RegisterConnection(conn)

	loggingHandler <- []byte("connection registred with id = 1.")

	// test the time out and the clean
	ticker := time.NewTicker(10 * time.Second)

	go func(t *time.Ticker) {
		// this will force the select to ignore
		// the ticker.C after 70seconds
		<-time.After(70 * time.Second)
		t.C = nil
	}(ticker)

	select {
	case tick := <-ticker.C:
		logMsg := fmt.Sprintf(
			"tick @ [%d:%d:%d] | total connections count: %d | cleaned connections : %d ",
			tick.Hour(),
			tick.Minute(),
			tick.Second(),
			connectionsPoolManager.GetConnectionsCount(),
			connectionsPoolManager.Clean(),
		)

		loggingHandler <- []byte(logMsg)

		// just for testing purposes, since the logger is not setup
		fmt.Println(logMsg)

	default:
		return
	}
}
