package networking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"subframe/server/settings"
	"subframe/server/storage"
	"subframe/structs/message"
)

var storageNodeActions = []string{
	"get",
	"put",
}

func startStorageNodeAPIService() {
	println("Starting HTTP Server at " + settings.StorageAddress)
	http.HandleFunc("/storage/", handleRequest)
	go func() {
		log.Fatal(http.ListenAndServe(settings.StorageAddress, nil))
	}()
}

func handleRequest(responseWriter http.ResponseWriter, req *http.Request) {
	println("Handling incoming request to " + req.URL.Path)
	request := storageRequest{
		res: responseWriter,
		req: req,
	}

	if request.parsePath() != http.StatusOK || !request.isValid() {
		writeResponse(responseWriter, http.StatusBadRequest, "Invalid Action or MessageID")
		return
	}

	//Handle Request
	request.handle()
}

type storageRequest struct {
	res       http.ResponseWriter
	req       *http.Request
	action    string
	messageID string
	valid     bool
}

func (r *storageRequest) parsePath() (status int) {
	parts := strings.Split(r.req.URL.Path, "/")
	if len(parts) < 3 {
		return http.StatusBadRequest
	}
	r.action = parts[2]
	r.messageID = regexp.MustCompile("[^A-Za-z0-9]").ReplaceAllString(parts[3], "-")
	return http.StatusOK
}

func (r *storageRequest) isValid() bool {
	validAction := false
	validMsgID := false
	for _, a := range storageNodeActions {
		if r.action == a {
			validAction = true
		}
	}
	if len(r.messageID) > 0 {
		validMsgID = true
	}
	r.valid = validAction && validMsgID
	return validAction && validMsgID
}

func (r storageRequest) handle() {
	//Handle request
	switch r.action {
	case "get":
		r.handleGet()
	case "put":
		r.handlePut()
	}
}

func (r storageRequest) handleGet() {
	message, readingError := storage.Get(r.messageID)
	if readingError != http.StatusOK {
		writeResponse(r.res, readingError, "Error getting message with ID "+r.messageID)
		return
	}
	responsedata, encodingError := json.Marshal(message)
	if encodingError != nil {
		writeResponse(r.res, http.StatusInternalServerError, "Error serving message from disk")
		return
	}
	writeResponse(r.res, http.StatusOK, string(responsedata))
}

func (r storageRequest) handlePut() {
	messageID := r.messageID
	r.req.Body = http.MaxBytesReader(r.res, r.req.Body, settings.MessageMaxSize*1024*1024)
	messageBody, error := ioutil.ReadAll(r.req.Body)
	if error != nil {
		if int64(len(messageBody)) >= settings.MessageMaxSize*1024*1024 {
			writeResponse(r.res, http.StatusRequestEntityTooLarge, "Message too large to be accepted by this node")
			return
		}
		writeResponse(r.res, http.StatusBadRequest, "Transmission of Message Body failed. Please try again.")
		return
	}

	message := message.Message{
		ID:      messageID,
		Content: string(messageBody),
	}

	status := storage.Put(message)

	if status != http.StatusOK {
		writeResponse(r.res, status, "Error storing message "+messageID)
		return
	}

	writeResponse(r.res, http.StatusOK, "Successfully stored message "+messageID)

	//Announce Message to CoordinatorNetwork
	//Push Message to other StorageNodes
}

func writeResponse(w http.ResponseWriter, status int, response string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, response)
}
