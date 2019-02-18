package main

import (
	"subframe/server/database"
	"subframe/server/settings"
)

var greeter = "                  .--.                  \n              `-/oooooo/-`              \n           .:+oooooooooooo+:.           \n       `-/oooooooooooooooooooo/-`       \n    .:+oooooooooooooooooooooooooo+:.    \n  :oooooooo+////////////////+oooooooo:  \n /ooooooo:`                  `:ooooooo/ \n /ooooooo   -//////////////-   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /oooooooooooooo/   ooooooo/ \n /oooooo+   /+-............`  `ooooooo/ \n /oooooo+   /oo:`           `-+ooooooo/ \n /oooooo+   /ooooooooooooooooooooooooo/ \n /oooooo+ `:oooooooooooooooooooooooooo/ \n  :oooooo:ooooooooooooooooooooooooooo:  \n    .:+oooooooooooooooooooooooooo+:.    \n       `-/oooooooooooooooooooo/-`       \n           .:+oooooooooooo+:.           \n              `-/oooooo/-`              \n                  .--.                  "

func main() {
	println(greeter)
	println("Initializing SuBFraMe Server...")
	settings.Read()
	//defer settings.Write()
	/*networking.Init()
	defer networking.Stop()
	storage.Init()
	defer storage.Finish()*/
	database.Init()
	defer database.Close()

	//jobqueue.SpawnWorker()

	/*for {
	}*/
}
