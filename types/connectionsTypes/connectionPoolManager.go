package connectionstypes

type ConnectionsPoolManager interface {
	// embbed the interface type to get
	// all the method set that defines it
	ConnectionsPool

	// extra functionnality
	GetLoggingHandler() chan<- []byte
	SetLoggingHandler(chan<- []byte)
}

type HttpConnectionPoolManager struct {
	ConnectionsPool

	// (NOTE) right now the logging handler does
	// not do much for HttpConnectionsPoolManager
	// but it will be useful if we had more methods
	// specific to this type where we can log diffrent
	// actions, like init, notifications...
	loggingHandler chan<- []byte
}

func (pm HttpConnectionPoolManager) GetLoggingHandler() chan<- []byte {
	return pm.loggingHandler
}
func (pm *HttpConnectionPoolManager) SetLoggingHandler(loggingHandler chan<- []byte) {
	pm.loggingHandler = loggingHandler
}
