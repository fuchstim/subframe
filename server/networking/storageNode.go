package networking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"subframe/server/database"
	"subframe/server/jobqueue"
	"subframe/server/logger"
	"subframe/server/settings"
	"subframe/server/storage"
	"subframe/structs/message"
)

var slog = logger.Logger{Prefix: "networking/StorageNode"}

var storageNodeActions = []string{
	"get",
	"put",
	"update",
	"control",
}

func startStorageNodeAPIService() {
	slog.Info("Starting HTTP Server at " + settings.LocalAddress + "...")
	http.HandleFunc("/storage/", handleRequest)
	go func() {
		err := (http.ListenAndServe(settings.LocalAddress, nil))
		slog.Fatal("Fatal failure in HTTP Storage Interface Server: " + err.Error())
	}()
}

func handleRequest(responseWriter http.ResponseWriter, req *http.Request) {
	slog.Info("Handling incoming " + req.Method + " request to " + req.URL.Path + "...")
	request := storageRequest{
		res: responseWriter,
		req: req,
	}

	if request.parsePath() != http.StatusOK || !request.isValid() {
		slog.Info("Action or Slug for " + req.URL.Path + " is invalid")
		writeResponse(responseWriter, http.StatusBadRequest, "Invalid Action or Slug")
		return
	}

	//Handle Request
	slog.Info("Request appears valid (Action: " + request.action + ", Slug: " + request.slug + "). Processing...")
	request.handle()
}

type storageRequest struct {
	res    http.ResponseWriter
	req    *http.Request
	action string
	slug   string
	valid  bool
}

func (r *storageRequest) parsePath() (status int) {
	parts := strings.Split(r.req.URL.Path, "/")[1:]
	if len(parts) < 2 {
		return http.StatusBadRequest
	} else if len(parts) < 3 {
		r.action = parts[1]
		return http.StatusOK
	}
	r.action = parts[1]
	rexp, err := regexp.Compile("[^A-Za-z0-9]")
	if err != nil {
		return http.StatusBadRequest
	}
	r.slug = rexp.ReplaceAllString(parts[2], "-")
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
	if len(r.slug) > 0 {
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
	case "control":
		r.handleControl()
	case "update":
		r.updateMessageStatus()
	}
}

func (r storageRequest) handleGet() {
	slog.Info("Handling MessageGET Request for " + r.slug + "...")

	if r.req.Method != "GET" {
		slog.Error("Client is trying to MessageGET with a " + r.req.Method + " Request.")
		writeResponse(r.res, http.StatusBadRequest, r.req.Method+" is not allowed here.")
		return
	}

	message, readingError := storage.Get(r.slug)
	if readingError != http.StatusOK {
		slog.Error("Cannot server Message " + r.slug + ": " + strconv.Itoa(readingError))
		writeResponse(r.res, readingError, "Error getting message with ID "+r.slug)
		return
	}
	responsedata, encodingError := json.Marshal(message)
	if encodingError != nil {
		slog.Error("Error serving Message " + r.slug + ": " + encodingError.Error())
		writeResponse(r.res, http.StatusInternalServerError, "Error serving message from disk")
		return
	}
	slog.Info("Serving Message " + r.slug + "...")
	writeResponse(r.res, http.StatusOK, string(responsedata))
}

func (r storageRequest) handlePut() {
	slog.Info("Handling MessagePUT Request for " + r.slug + "...")

	if r.req.Method != "POST" {
		slog.Error("Client is trying to MessagePUT with a " + r.req.Method + " Request.")
		writeResponse(r.res, http.StatusBadRequest, r.req.Method+" is not allowed here.")
		return
	}

	messageID := r.slug
	r.req.Body = http.MaxBytesReader(r.res, r.req.Body, int64(settings.MessageMaxSize)*1024*1024)
	messageBody, error := ioutil.ReadAll(r.req.Body)
	if error != nil {
		if len(messageBody) >= settings.MessageMaxSize*1024*1024 {
			exceeds := (len(messageBody) / 1024 / 1024) - settings.MessageMaxSize
			slog.Error("Message size exceeds settings.MessageMaxSize (by " + strconv.Itoa(exceeds) + "M), denying storage request.")
			writeResponse(r.res, http.StatusRequestEntityTooLarge, "Message too large to be accepted by this node")
			return
		}
		slog.Error("Transmission of message failed: " + error.Error())
		writeResponse(r.res, http.StatusBadRequest, "Transmission of Message Body failed. Please try again.")
		return
	}

	//TODO: Verify that message is somewhat valid
	if len(messageBody) == 0 {
		slog.Error("Message Body is empty")
		writeResponse(r.res, http.StatusBadRequest, "Empty Message Body")
		return
	}

	slog.Info("Message " + messageID + " successfully transmitted. Storing...")
	message := message.Message{
		ID:      messageID,
		Content: string(messageBody),
	}

	status := storage.Put(message)
	if status == http.StatusOK {
		status = database.LogMessageStorage(messageID)
	}

	if status != http.StatusOK {
		slog.Error("Error storing message: " + strconv.Itoa(status))
		writeResponse(r.res, status, "Error storing message "+messageID)
		return
	}

	slog.Info("Successfully stored Message " + messageID)
	writeResponse(r.res, http.StatusOK, "Successfully stored message "+messageID)

	task := func(data interface{}) {
		log := logger.Logger{Prefix: "networking/Announce-" + messageID}
		messageID, ok := data.(string)
		if !ok {
			log.Error("Error starting Announcing Thread")
			return
		}

		log.Info("Getting CoordinatorNodes to announce Message to...")
		//Get three random coordinatorNodes
		coordinatorNodes := database.GetRandomCoordinatorNodes(3)
		if len(coordinatorNodes) == 0 {
			log.Error("Received empty List of CoordinatorNodes.")
			return
		}
		log.Info("Announcing Message to " + strconv.Itoa(len(coordinatorNodes)) + " CoordinatorNodes...")
		//Announce MessageID to CoordinatorNetwork
		var redistribute = "true"
		for _, value := range coordinatorNodes {
			r := SendNodeRequest(NODE_COORDINATOR, value, "/announce/"+messageID+"/"+settings.RemoteAddress, "")
			//If at least one node orders to not further distribute the message, do not
			if r == "false" {
				redistribute = r
			}
		}
		log.Info("Announced Message to CoordinatorNetwork. Redistributing: " + redistribute)
		if redistribute == "true" {
			//TODO: Push Message to other StorageNodes
		}
	}
	job := jobqueue.Job{
		Task: task,
		Data: messageID,
	}
	select {
	case jobqueue.Queue <- job:
	}
}

func (r storageRequest) handleControl() {
	action := r.slug
	switch action {
	case "get-storage-nodes":
		r.printStorageNodes()
	case "get-coordinator-nodes":
		r.printCoordinatorNodes()
	}
}

func (r storageRequest) printStorageNodes() {
	slog.Info("Exporting 10 StorageNodes...")
	storageNodes := database.GetStorageNodes(10)
	response, err := json.Marshal(storageNodes)
	if err != nil {
		slog.Error("Failed to export StorageNodes: " + err.Error())
		writeResponse(r.res, http.StatusInternalServerError, "Failed to export StorageNodes.")
		return
	}
	slog.Info("Exported StorageNodes.")
	writeResponse(r.res, http.StatusOK, string(response))
}

func (r storageRequest) printCoordinatorNodes() {
	slog.Info("Exporting CoordinatorNodes...")
	coordinatorNodes := database.GetCoordinatorNodes()
	response, err := json.Marshal(coordinatorNodes)
	if err != nil {
		slog.Error("Failed to export CoordinatorNodes: " + err.Error())
		writeResponse(r.res, http.StatusInternalServerError, "Failed to export CoordinatorNodes.")
		return
	}
	slog.Info("Exported CoordinatorNodes.")
	writeResponse(r.res, http.StatusOK, string(response))
}

func (r storageRequest) updateMessageStatus() {
	slog.Info("Received UPDATE for Message " + r.slug)
	messageID := r.slug

	job := jobqueue.Job{
		Task: func(data interface{}) {
			messageID, ok := data.(string)
			if ok {
				log := logger.Logger{Prefix: "networking/Update-" + messageID}
				status := GetMessageStatus(messageID)
				if status > -1 {
					log.Info("Updating Message Status to " + strconv.Itoa(status))
					database.UpdateMessageStatusStorage(messageID, status)
					return
				}
				log.Error("Received inconclusive Message Status. Not updating local database.")
			} else {
				slog.Error("Error Starting Update-Thread")
			}
		},
		Data: messageID,
	}

	select {
	case jobqueue.Queue <- job:
	}

	writeResponse(r.res, http.StatusOK, "OK")
}

func writeResponse(w http.ResponseWriter, status int, response string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, response)
}
