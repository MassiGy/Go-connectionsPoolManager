package main

import (
	"connectionsPoolManager/types"
	"context"
	"fmt"
	"sync"
	"time"
)

var connectionsLimitRate int
var connection types.Connection
var connectionsPool types.ConnectionsPool

/*
	we can not declare the connectionsPoolManager as a variable
	that implements ConnectionsPoolManager interface since we
	need to acces the embbeded type ConnectionsPool, and we
	can only do that through a concret type. Nevertheless,
	connectionsPoolManager variables can be considered as
	a types.ConnectionsPoolManager statifyign object, so we can
	pass it to any function that requires that interface

	(Note)	This was what causing the seg fault

	var connectionsPoolManager types.ConnectionsPoolManager
*/

var loggingHandler chan []byte
var waitGroup sync.WaitGroup

func main() {
	// set the limit of connection
	connectionsLimitRate = 8

	// init a connections pool
	connectionsPool = &types.HttpConnectionPool{
		Connections: make([]types.HttpConnection, 0, connectionsLimitRate),
	}
	connectionsPool.SetPoolSize(connectionsLimitRate) // pool size should'nt be exposed

	// init a connection pool manager
	connectionsPoolManager := &types.HttpConnectionPoolManager{}
	connectionsPoolManager.ConnectionsPool = connectionsPool

	loggingHandler = make(chan []byte) // @todo: setup an actual logger

	waitGroup.Add(1) // register following go routine
	go func() {
		// consume the logs to prevent blocking
		// as long as the channel is not closed
		for {
			msg, ok := <-loggingHandler

			if !ok {
				// break to unregister self from waitGrp
				break
			}

			// print the message
			fmt.Println(string(msg))
		}

		waitGroup.Done()
	}()

	connectionsPoolManager.SetLoggingHandler(loggingHandler)

	loggingHandler <- []byte("initialization of connection pools done.")

	loggingHandler <- []byte("about to create a connection.")

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))

	connection = &types.HttpConnection{}
	connection.SetId(1)
	connection.SetContext(ctx)

	connectionsPoolManager.RegisterConnection(connection)

	loggingHandler <- []byte("connection registred with id = 1.")

	// test the time out and the clean
	ticker := time.NewTicker(5 * time.Second)
	timeOut := make(chan int)

	waitGroup.Add(1) // register following go routine
	go func() {
		// this will force the select to ignore
		// the ticker.C after 70seconds and enter
		// the second case to quit the program
		<-time.After(70 * time.Second)
		timeOut <- 1

		ticker.C = nil

		waitGroup.Done()
	}()

	for {

		select {
		case tick := <-ticker.C:
			logMsg := fmt.Sprintf(
				"tick @ [%d:%d:%d]\t|total connections count: %d\t|cleaned connections : %d",
				tick.Hour(),
				tick.Minute(),
				tick.Second(),
				connectionsPoolManager.ConnectionsPool.GetConnectionsCount(),
				connectionsPoolManager.ConnectionsPool.Clean(),
			)

			loggingHandler <- []byte(logMsg)

		case <-timeOut:
			// stop the logs consumption go routine
			close(loggingHandler)

			// wait for all goroutines to finish
			waitGroup.Wait()

			// stop the demo and quit
			fmt.Println("quitting...")
			return
		}
	}
}
