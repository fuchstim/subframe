package main

import (
	"subframe/server/jobqueue"
	"subframe/server/networking"
	"subframe/server/settings"
	"subframe/server/storage"
)

func main() {
	println("Initializing Server...")
	settings.Read()
	networking.Init()
	storage.Init()

	jobqueue.SpawnWorker()

	for {
	}
}
