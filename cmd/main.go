package main

import (
	"connectionsPoolManager/types"
	"context"
	"fmt"
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

// var loggingHandler chan []byte

func main() {
	// set the limit of connection
	connectionsLimitRate = 8

	// init a connections pool
	connectionsPool = &types.HttpConnectionPool{
		Connections: make([]types.HttpConnection, 0, connectionsLimitRate),
	}
	connectionsPool.SetPoolSize(connectionsLimitRate)

	// init a connection pool manager
	connectionsPoolManager := &types.HttpConnectionPoolManager{}
	connectionsPoolManager.ConnectionsPool = connectionsPool

	/*
		loggingHandler = make(chan []byte) // @todo: setup an actual logger
		go func() {
			// consume the logs to prevent blocking
			fmt.Println(string(<-loggingHandler))
		}()

		connectionsPoolManager.SetLoggingHandler(loggingHandler)

		loggingHandler <- []byte("initialization of connection pools done.")

		loggingHandler <- []byte("about to create a connection.")
	*/
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))

	connection = &types.HttpConnection{}
	connection.SetId(1)
	connection.SetContext(ctx)

	connectionsPoolManager.RegisterConnection(connection)

	//loggingHandler <- []byte("connection registred with id = 1.")

	// test the time out and the clean
	ticker := time.NewTicker(time.Second)
	timeOut := make(chan int)

	go func() {
		// this will force the select to ignore
		// the ticker.C after 70seconds and enter
		// the second case to quit the program
		<-time.After(70 * time.Second)
		timeOut <- 1

		ticker.C = nil
	}()

	for {

		select {
		case tick := <-ticker.C:
			logMsg := fmt.Sprintf(
				"tick @ [%d:%d:%d] | total connections count: %d | cleaned connections : %d ",
				tick.Hour(),
				tick.Minute(),
				tick.Second(),
				connectionsPoolManager.ConnectionsPool.GetConnectionsCount(),
				connectionsPoolManager.ConnectionsPool.Clean(),
			)

			//loggingHandler <- []byte(logMsg)

			// just for testing purposes, since the logger is not setup
			fmt.Println(string(logMsg))

		/*

			We should not set time.After in here since it
			will be reset to 0s each iteration, so this will
			take more then 2 minutes to kick in.
		*/
		case <-time.After(2 * time.Minute):
			// @todo: stop all the conccurent go routines
			// stop the demo and quit
			fmt.Println("returning...")
			return

		}
	}
}
