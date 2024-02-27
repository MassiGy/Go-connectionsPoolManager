package connectionstypes

import "errors"

type ConnectionsPool interface {
	GetPoolSize() int // > 0, if 0 then non-initialized
	SetPoolSize(int)
	GetConnectionsCount() int

	RegisterConnection(Connection) error
	KillConnection(Connection)
	Clean() int // returns how many non-alive connections cleaned
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
