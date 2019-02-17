package networking

//Init Initializes StorageNode HTTP Api and starts coordinator network service
func Init() {
	println("Initializing Networking...")
	//Start StorageNode Api
	startStorageNodeAPIService()

	//Start CoordinatorNode service

}
