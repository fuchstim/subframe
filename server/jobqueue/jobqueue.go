package jobqueue

import (
	"subframe/server/settings"
)

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
	die chan bool
}

func (sw worker) start() {
	go func() {
		for {
			workerCount := len(workerPool)
			queueLength := len(Queue)

			if workerCount < settings.MaxWorkers && queueLength >= settings.QueueMaxLength {
				println("Spawning new Worker...")
				SpawnWorker()
			} else if workerCount > 1 && queueLength <= settings.QueueMaxLength {
				println("Killing Worker")
				sw.die <- true
			}
			select {
			case job := <-Queue:
				{
					println("Executing job...")
					job.execute()
				}
			case <-sw.die:
				{
					println("Stopped worker")
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
		return
	}
	println("Spawning and starting new Worker...")
	worker := worker{
		die: make(chan bool),
	}
	workerPool = append(workerPool, &worker)
	worker.start()
}
