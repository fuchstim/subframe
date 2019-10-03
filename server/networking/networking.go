package networking

import (
	"subframe/server/logger"
	. "subframe/status"
)

var mlog = logger.Logger{Prefix: "networking/Main"}

//Init Initializes StorageNode HTTP Api and starts coordinator network service
func Init() {
	mlog.Info(InProgress, "Initializing Networking...")
	//Start StorageNode Api
	startStorageNodeAPIService()

	//Start CoordinatorNode service
	mlog.Info(OK, "Initialized Networking.")
}

//Stop terminates and stops all active network connections and interfaces
func Stop() {
	mlog.Info(InProgress, "Stopping Networking...")

	mlog.Info(OK, "Stopped Networking.")
}
