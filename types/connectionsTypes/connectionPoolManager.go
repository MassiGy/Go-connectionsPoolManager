package connectionstypes

import (
	loggertypes "connectionsPoolManager/types/loggerTypes"
	"fmt"
	"time"
)

type ConnectionsPoolManager interface {
	// embbed the interface type to get
	// all the method set that defines it
	ConnectionsPool

	// extra functionnality
	GetLoggingHandler() loggertypes.LogQueuer
	SetLoggingHandler(loggertypes.LogQueuer)
	Monitor(monitoringPeriod time.Duration, monitoringIntervals time.Duration)
}

type HttpConnectionPoolManager struct {
	ConnectionsPool

	loggingHandler loggertypes.LogQueuer
}

func (pm HttpConnectionPoolManager) GetLoggingHandler() loggertypes.LogQueuer {
	return pm.loggingHandler
}
func (pm *HttpConnectionPoolManager) SetLoggingHandler(loggingHandler loggertypes.LogQueuer) {
	pm.loggingHandler = loggingHandler
}

func (pm *HttpConnectionPoolManager) Monitor(monitoringPeriod time.Duration, monitoringIntervals time.Duration) {
	ticker := time.NewTicker(monitoringIntervals)
	timeoutChan := time.After(monitoringPeriod)

	for {

		select {
		case <-ticker.C:
			logMsg := fmt.Sprintf(
				"Total connections count: %d\t|cleaned connections : %d\n",
				pm.ConnectionsPool.GetConnectionsCount(),
				pm.ConnectionsPool.Clean(),
			)
			pm.loggingHandler.EnQueue([]byte(logMsg))

		case <-timeoutChan:
			pm.loggingHandler.EnQueue([]byte("Ending monitoring for the connection pool manager\n"))
			return
		}
	}
}
