package settings

import (
	"flag"
)

//BootstrapNode is used for Bootstrapping the local instance
var BootstrapNode string = ""

//DataPath is used to store message files
var DataPath string = "./data"

//InterfaceAddress is used to access the local instance remotely
var InterfaceAddress string = "localhost:9123"

//StorageAddress is the IP and Port the StorageNode instance listens on
var StorageAddress string = "0.0.0.0:9123"

//DiskSpace is the maximum space used for message storage
var DiskSpace int = 5000

//MaxWorkers is the maximum number of workers to spawn
var MaxWorkers int = 10

//QueueMaxLength is the maximum length of a queue before a new worker is spawned, if the current worker count does not exceed MaxWorkers
var QueueMaxLength int = 10

//MessageMaxSize defines the maximum size of an individual message file
var MessageMaxSize int64 = 100

//MessageMinCheckDelay defines the minimum time in hours between individual checks of the message status
var MessageMinCheckDelay int = 12

//MessageMaxStoreTime defines the maximum time a message is stored locally, in days
var MessageMaxStoreTime int = 7

//Read reads settings from local storage and overwrites them with command-line-arguments
func Read() {
	println("Reading Settings...")
	///Read Settings from disk, then overwrite with command line args

	parseCommandLineArgs()
	println("Read Settings.")
}

//Write writes settings to local storage
func Write(path string, settings *map[string]string) {
	//Write settings to disk
}

func parseCommandLineArgs() {
	flag.StringVar(&BootstrapNode, "bootstrap-node", BootstrapNode, "If set, SuBFraMe will reinitialize the local Node Database and sync it with the BootstrapNode")
	flag.StringVar(&DataPath, "data-dir", DataPath, "The SuBFraMe data directory, messages will be stored here")
	flag.StringVar(&InterfaceAddress, "interface-addr", InterfaceAddress, "The address of this SuBFraMe Instance")
	flag.StringVar(&StorageAddress, "storage-address", StorageAddress, "The IP and Port the StorageNode Interface will listen on")
	flag.IntVar(&DiskSpace, "disk-space", DiskSpace, "The maximum space SuBFraMe will use to store Messages in MB")
	flag.IntVar(&MaxWorkers, "max-workers", MaxWorkers, "The maximum number of worker threads")
	flag.IntVar(&QueueMaxLength, "max-queue-length", QueueMaxLength, "The maximum size a queue can have before a new worker is spawned, before exceeding max-workers")
	flag.Int64Var(&MessageMaxSize, "message-max-size", MessageMaxSize, "The maximum size of an individual message file, in MB")
	flag.IntVar(&MessageMinCheckDelay, "message-min-check-delay", MessageMinCheckDelay, "The minimum time in hours between individual checks of the same message against the coordinator network")
	flag.IntVar(&MessageMaxStoreTime, "message-max-store-time", MessageMaxStoreTime, "The maximum time a message is stored locally, in days")
	flag.Parse()
}
