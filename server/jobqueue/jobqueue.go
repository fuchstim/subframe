package jobqueue

import (
	"strconv"
	"subframe/server/logger"
	"subframe/server/settings"
	"time"
)

var log = logger.Logger{Prefix: "jobqueue/Main"}

//Task will be executed by Job
type Task func(data interface{})

//Job will be executed
type Job struct {
	Task Task
	Data interface{}
}

func (j Job) execute() {
	j.Task(j.Data)
}

type worker struct {
	id  string
	die chan bool
}

func (sw worker) start() {
	log.Info("Starting worker " + sw.id + ".")
	go func() {
		for {
			workerCount := len(workerPool)
			queueLength := len(Queue)

			if workerCount < settings.MaxWorkers && queueLength >= settings.QueueMaxLength {
				log.Info("Queue length exceeds settings.MaxQueueLength.")
				SpawnWorker()
			} else if workerCount > 1 && queueLength <= settings.QueueMaxLength {
				log.Info("Too many workers for current queue length. Killing worker " + sw.id + "...")
				sw.die <- true
			}
			select {
			case job := <-Queue:
				{
					job.execute()
				}
			case <-sw.die:
				{
					log.Info("Worker " + sw.id + " killed.")
					removeWorkerFromPool(sw.id)
					return
				}
			}
		}
	}()
}

var workerPool []*worker

//Queue holds all jobs waiting to be executed
var Queue = make(chan Job)

//SpawnWorker spawns a new Worker, if MaxWorkers setting allows it
func SpawnWorker() {
	if len(workerPool) >= settings.MaxWorkers {
		log.Warn("settings.MaxWorkers does not allow for a new Worker to be spawned.")
		return
	}
	log.Info("Spawning and starting new Worker...")
	worker := worker{
		id:  strconv.FormatInt(time.Now().Unix(), 16),
		die: make(chan bool),
	}
	workerPool = append(workerPool, &worker)
	log.Info("New worker count: " + strconv.Itoa(len(workerPool)))
	worker.start()
}

func removeWorkerFromPool(id string) {
	//Remove worker with id from pool
	for index, value := range workerPool {
		if value.id == id {
			workerPool = append(workerPool[:index], workerPool[index+1:]...)
		}
	}
}
