package clientdb

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ryrden/rinha-de-backend-go/internal/app/domain/client"
)

var (
	MaxWorker = 8
	MaxQueue  = 3
)

type JobQueue chan Job

// Job represents the job to be run
type Job struct {
	Payload *ClientTransactionPayload
}

type ClientTransactionPayload struct {
	Client      *client.Client
	Value       int
	Kind        string
	Description string
}

// A buffered channel that we can send work requests on.

func NewJobQueue() JobQueue {
	return make(JobQueue, MaxQueue)
}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
	dbPool     *pgxpool.Pool
}

func (w Worker) Start() {
	dataCh := make(chan Job)
	insertCh := make(chan []Job)

	go w.bootstrap(dataCh)

	go w.processData(dataCh, insertCh)

	go w.processInsert(insertCh)
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func (w Worker) bootstrap(dataCh chan Job) {
	for {
		w.WorkerPool <- w.JobChannel

		select {
		case job := <-w.JobChannel:
			dataCh <- job

		case <-w.quit:
			return
		}
	}
}

func (w Worker) processData(dataCh chan Job, insertCh chan []Job) {
	tickInsertRate := time.Duration(10 * time.Second)
	tickInsert := time.NewTicker(tickInsertRate).C

	batchMaxSize := 5000
	batch := make([]Job, 0, batchMaxSize)

	for {
		select {
		case data := <-dataCh:
			batch = append(batch, data)

		case <-tickInsert:
			batchLen := len(batch)
			if batchLen > 0 {
				log.Infof("Tick insert (len=%d)", batchLen)
				insertCh <- batch

				batch = make([]Job, 0, batchMaxSize)
			}
		}
	}
}

func (w Worker) processInsert(insertCh chan []Job) {
	for {
		select {
		case jobs := <-insertCh:
			for _, job := range jobs {
				transactionData := job.Payload
				_, err := w.dbPool.Exec(
					context.Background(),
					"INSERT INTO transactions(client_id, amount, kind, description) VALUES($1, $2, $3, $4)",
					transactionData.Client.ID,
					transactionData.Value,
					transactionData.Kind,
					transactionData.Description,
				)
				if err != nil {
					log.Errorf("Error on insert transaction: %v", err)
				}
				log.Infof("Transaction inserted successfully")
			}
		case <-w.quit:
			return
		}
	}
}

func NewWorker(workerPool chan chan Job, db *pgxpool.Pool) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		dbPool:     db,
	}
}

type Dispatcher struct {
	maxWorkers int
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	jobQueue   chan Job
	dbPool     *pgxpool.Pool
}

func NewDispatcher(db *pgxpool.Pool, jobQueue JobQueue) *Dispatcher {
	maxWorkers := MaxWorker

	pool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
		jobQueue:   jobQueue,
		dbPool:     db,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.dbPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
