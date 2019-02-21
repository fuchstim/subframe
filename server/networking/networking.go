package networking

import (
	"subframe/server/logger"
)

var mlog = logger.Logger{Prefix: "networking/Main"}

//Init Initializes StorageNode HTTP Api and starts coordinator network service
func Init() {
	mlog.Info("Initializing Networking...")
	//Start StorageNode Api
	startStorageNodeAPIService()

	//Start CoordinatorNode service
	mlog.Info("Initialized Networking.")
}

//Stop terminates and stops all active network connections and interfaces
func Stop() {
	mlog.Info("Stopping Networking...")

	mlog.Info("Stopped Networking.")
}
