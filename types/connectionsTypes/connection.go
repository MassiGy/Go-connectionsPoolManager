package connectionstypes

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
