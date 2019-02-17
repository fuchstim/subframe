package settings

import (
	"flag"
)

//BootstrapNode is used for Bootstrapping the local instance
var BootstrapNode string

//DataPath is used to store message files
var DataPath string

//InterfaceAddress is used to access the local instance remotely
var InterfaceAddress string

//StorageAddress is the IP and Port the StorageNode instance listens on
var StorageAddress string

//DiskSpace is the maximum space used for message storage
var DiskSpace int

//MaxWorkers is the maximum number of workers to spawn
var MaxWorkers int

//MaxQueueLength is the maximum length of a queue before a new worker is spawned, if the current worker count does not exceed MaxWorkers
var MaxQueueLength int

//MaxMessageSize defines the maximum size of an individual message file
var MaxMessageSize int64

//Read reads settings from local storage and overwrites them with command-line-arguments
func Read() {
	println("Reading Settings...")
	///Read Settings from disk, then overwrite with command line args

	parseCommandLineArgs()
	println("Read Settings.")
}

func Write(path string, settings *map[string]string) {
	//Write settings to disk
}

func parseCommandLineArgs() {
	flag.StringVar(&BootstrapNode, "bootstrap-node", "", "If set, SuBFraMe will reinitialize the local Node Database and sync it with the BootstrapNode")
	flag.StringVar(&DataPath, "data-dir", "./data", "The SuBFraMe data directory, messages will be stored here")
	flag.StringVar(&InterfaceAddress, "interface-addr", "localhost:9123", "The address of this SuBFraMe Instance")
	flag.StringVar(&StorageAddress, "storage-address", "0.0.0.0:9123", "The IP and Port the StorageNode Interface will listen on")
	flag.IntVar(&DiskSpace, "disk-space", 5000, "The maximum space SuBFraMe will use to store Messages in MB")
	flag.IntVar(&MaxWorkers, "max-workers", 10, "The maximum number of worker threads")
	flag.IntVar(&MaxQueueLength, "max-queue-length", 10, "The maximum size a queue can have before a new worker is spawned, before exceeding max-workers")
	flag.Int64Var(&MaxMessageSize, "max-message-size", 100, "The maximum size of an individual message file, in MB")
	flag.Parse()
}
