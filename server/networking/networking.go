package networking

//Init Initializes StorageNode HTTP Api and starts coordinator network service
func Init() {
	println("Initializing Networking...")
	//Start StorageNode Api
	startStorageNodeAPIService()

	//Start CoordinatorNode service

}

//Stop terminates and stops all active network connections and interfaces
func Stop() {

}
