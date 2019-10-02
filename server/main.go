package main

import (
	"os"
	"os/signal"
	"subframe/server/bootstrapper"
	"subframe/server/database"
	"subframe/server/jobqueue"
	"subframe/server/logger"
	"subframe/server/networking"
	"subframe/server/settings"
	"subframe/server/storage"
)

var log = logger.Logger{Prefix: "main/Main"}
var greeter = "                  .--.                  \n              `-/oooooo/-`              \n           .:+oooooooooooo+:.           \n       `-/oooooooooooooooooooo/-`       \n    .:+oooooooooooooooooooooooooo+:.    \n  :oooooooo+////////////////+oooooooo:  \n /ooooooo:`                  `:ooooooo/ \n /ooooooo   -//////////////-   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /+-............`  `ooooooo/ \n /oooooo+   /oo:`           `-+ooooooo/ \n /oooooo+   /ooooooooooooooooooooooooo/ \n /oooooo+ `:oooooooooooooooooooooooooo/ \n  :oooooo:ooooooooooooooooooooooooooo:  \n    .:+oooooooooooooooooooooooooo+:.    \n       `-/oooooooooooooooooooo/-`       \n           .:+oooooooooooo+:.           \n              `-/oooooo/-`              \n                  .--.                  "

func main() {
	println(greeter)
	log.Info("Welcome to SuBFraMe Server!")
	log.Info("Initializing Server...")

	settings.Read()
	defer settings.Write()

	storage.Init()
	defer storage.Finish()

	logger.Init()
	defer logger.Close()

	jobqueue.SpawnWorker()

	networking.Init()
	defer networking.Stop()

	database.Init()
	defer database.Close()

	bootstrapper.Bootstrap()

	//Wait for interrupt, then return
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Info("Stopping SuBFraMe Server...")
	return
}
