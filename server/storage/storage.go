package storage

import (
	"subframe/structs/message"
)

func Get(id string) (msg message.Message, status int) {
	//Read message from disk and return
	return message.Message{}, 0
}

func Put(msg message.Message) (status int) {
	//Write message to disk
	return 0
}

func Delete(id string) (status int) {
	//Delete message from disk
	return 0
}
