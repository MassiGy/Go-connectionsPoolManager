package main

import (
	connTypes "connectionsPoolManager/types/connectionsTypes"
	"context"
	"fmt"
	"sync"
	"time"
)

var connectionsLimitRate int
var connection connTypes.Connection
var connectionsPool connTypes.ConnectionsPool

// for connectionsPoolManager we need a concret type to
// access the embbeded type connTypes.ConnectionPool, we can
// not do so through an interface based type declaration
var connectionsPoolManager *connTypes.HttpConnectionPoolManager

var loggingHandler chan []byte //@todo: setup an actual logger

// wGrp to keep track of our goroutines
var waitGroup sync.WaitGroup

func main() {
	loggingHandler = make(chan []byte)

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

	connectionsLimitRate = 8
	loggingHandler <- []byte(fmt.Sprintf("Connections pool size set to %d.", connectionsLimitRate))

	// init a connections pool
	connectionsPool = &connTypes.HttpConnectionPool{
		Connections: make([]connTypes.HttpConnection, 0, connectionsLimitRate),
	}
	connectionsPool.SetPoolSize(connectionsLimitRate) // pool size should'nt be exposed

	// init a connection pool manager
	connectionsPoolManager = &connTypes.HttpConnectionPoolManager{}
	connectionsPoolManager.ConnectionsPool = connectionsPool

	/*
		Right now we do not have any methods in the Connections
		Pool Manager interface that take advantage of the logger

		connectionsPoolManager.SetLoggingHandler(loggingHandler)
	*/
	loggingHandler <- []byte("Initialization of connection pools done.")

	loggingHandler <- []byte("About to create a connection.")

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))

	connection = &connTypes.HttpConnection{}
	connection.SetId(1)
	connection.SetContext(ctx)

	connectionsPoolManager.RegisterConnection(connection)

	loggingHandler <- []byte(fmt.Sprintf("Connection registred with id = %d.", connection.GetId()))

	// test the time out and the clean
	ticker := time.NewTicker(5 * time.Second)
	loggingHandler <- []byte("Connections Pool monitoring time cycle set to 5s.")

	timeOut := make(chan int)

	waitGroup.Add(1) // register following go routine
	go func() {
		// this will force the select to ignore
		// the ticker.C after 70seconds and enter
		// the second case to quit the program
		loggingHandler <- []byte("Connections Pool monitoring time out set to 70s.")

		<-time.After(70 * time.Second)
		timeOut <- 1

		ticker.C = nil

		waitGroup.Done()
	}()

	for {

		select {
		case tick := <-ticker.C:
			logMsg := fmt.Sprintf(
				"Tick @ [%d:%d:%d]\t|total connections count: %d\t|cleaned connections : %d",
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
			fmt.Println("Quitting...")
			return
		}
	}
}
