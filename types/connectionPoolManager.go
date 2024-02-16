package types

import (
	"context"
	"errors"
	"time"
)

type connectionType int

const (
	HTTP connectionType = iota // starts at 0
	GRPC
	WS
)

type Connection interface {
	GetId() int
	SetId(int) error
	GetConnectionType() connectionType
	IsAlive() bool
	GetContext() context.Context
	SetContext(ctx context.Context)
}

type ConnectionsPool interface {
	GetPoolSize() int // > 0, if 0 then non-initialized
	SetPoolSize(int)
	GetConnectionsCount() int

	RegisterConnection(Connection) error
	KillConnection(Connection)
	Clean() int // returns how many non-alive connections cleaned
}

type ConnectionsPoolManager interface {
	// embbed the interface type to get
	// all the method set that defines it
	ConnectionsPool

	// extra functionnality
	GetLoggingHandler() chan<- []byte
	SetLoggingHandler(chan<- []byte)
}

type HttpConnection struct {
	id int // > 0, if 0 then non-initialized

	// use context to manage the lifecycles
	ctx context.Context
}

func (c HttpConnection) GetId() int {
	return c.id
}
func (c *HttpConnection) SetId(id int) error {
	if id == 0 || c.id > 0 {
		return errors.New("Can not set id")
	}
	c.id = id
	return nil
}
func (c HttpConnection) GetConnectionType() connectionType {
	return HTTP
}

func (c HttpConnection) GetContext() context.Context {
	return c.ctx
}
func (c *HttpConnection) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c HttpConnection) IsAlive() bool {

	deadline, isset := c.ctx.Deadline()

	// if !isset, then no deadline is set
	return !isset || time.Now().Before(deadline)
}

type HttpConnectionPool struct {
	poolSize    int
	Connections []HttpConnection
}

func (p HttpConnectionPool) GetPoolSize() int {
	return p.poolSize
}
func (p *HttpConnectionPool) SetPoolSize(size int) {
	p.poolSize = size
}

func (p HttpConnectionPool) GetConnectionsCount() int {
	return len(p.Connections)
}

func (p *HttpConnectionPool) RegisterConnection(c Connection) error {
	if c == nil || !c.IsAlive() {
		return errors.New("Connection not valid")
	}

	if cap(p.Connections) < 1 {
		return errors.New("Pool is full.")
	}

	// type assertion (down cast c to HttpConnection)
	httpCon := c.(*HttpConnection)

	p.Connections = append(p.Connections, *httpCon)
	return nil
}

func (p *HttpConnectionPool) KillConnection(c Connection) {

	for i, conn := range p.Connections {
		if c.GetId() == conn.GetId() {

			// @todo kill the connection

			// safe removal due to return
			p.Connections = append(p.Connections[:i], p.Connections[i+1:]...)
		}
	}
}

func (p *HttpConnectionPool) Clean() int {
	cleaned := 0

	for _, conn := range p.Connections {

		// if id>0, conn is not initialized
		if conn.id != 0 && !conn.IsAlive() {
			cleaned++

			// these fn calls will be fired in
			// LIFO order at the return point
			defer p.KillConnection(&conn)
		}
	}

	return cleaned
}

type HttpConnectionPoolManager struct {
	ConnectionsPool
	loggingHandler chan<- []byte
}

func (pm HttpConnectionPoolManager) GetLoggingHandler() chan<- []byte {
	return pm.loggingHandler
}
func (pm *HttpConnectionPoolManager) SetLoggingHandler(loggingHandler chan<- []byte) {
	pm.loggingHandler = loggingHandler
}
