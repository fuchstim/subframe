package main

import (
	"subframe/server/database"
	"subframe/server/jobqueue"
	"subframe/server/networking"
	"subframe/server/settings"
	"subframe/server/storage"
)

func main() {
	println("Initializing Server...")
	settings.Read()
	//defer settings.Write()
	networking.Init()
	defer networking.Stop()
	storage.Init()
	defer storage.Finish()
	database.Init()
	defer database.Close()

	jobqueue.SpawnWorker()

	for {
	}
}
