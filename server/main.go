package main

import (
	"subframe/server/bootstrapper"
	"subframe/server/database"
	"subframe/server/jobqueue"
	"subframe/server/networking"
	"subframe/server/settings"
	"subframe/server/storage"
)

var greeter = "                  .--.                  \n              `-/oooooo/-`              \n           .:+oooooooooooo+:.           \n       `-/oooooooooooooooooooo/-`       \n    .:+oooooooooooooooooooooooooo+:.    \n  :oooooooo+////////////////+oooooooo:  \n /ooooooo:`                  `:ooooooo/ \n /ooooooo   -//////////////-   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /+-............`  `ooooooo/ \n /oooooo+   /oo:`           `-+ooooooo/ \n /oooooo+   /ooooooooooooooooooooooooo/ \n /oooooo+ `:oooooooooooooooooooooooooo/ \n  :oooooo:ooooooooooooooooooooooooooo:  \n    .:+oooooooooooooooooooooooooo+:.    \n       `-/oooooooooooooooooooo/-`       \n           .:+oooooooooooo+:.           \n              `-/oooooo/-`              \n                  .--.                  "

func main() {
	println(greeter)
	println("Initializing SuBFraMe Server...")
	jobqueue.SpawnWorker()
	settings.Read()
	//defer settings.Write()
	networking.Init()
	defer networking.Stop()
	storage.Init()
	defer storage.Finish()
	database.Init()
	defer database.Close()
	bootstrapper.Bootstrap()

	for {
	}
}
