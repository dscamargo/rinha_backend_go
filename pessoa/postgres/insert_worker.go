package postgres

import (
	"context"
	"github.com/dscamargo/rinha_backend_go/pessoa"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	MaxQueue   = 1
	MaxWorkers = 1
)

type JobQueue chan Job

type Job struct {
	Payload *pessoa.Pessoa
}

type Dispatcher struct {
	maxWorkers int
	WorkerPool chan chan Job
	jobQueue   chan Job
	db         *pgxpool.Pool
}

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
	db         *pgxpool.Pool
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
	tickInsertRate := 3 * time.Second
	tickInsert := time.Tick(tickInsertRate)

	batchMaxSize := 10000
	batch := make([]Job, 0, batchMaxSize)

	for {
		select {
		case data := <-dataCh:
			batch = append(batch, data)

		case <-tickInsert:
			if len(batch) > 0 {
				log.Infof("Insert - Size: %d", len(batch))
				insertCh <- batch

				batch = make([]Job, 0, batchMaxSize)
			}
		}
	}
}

func (w Worker) processInsert(insertCh chan []Job) {
	columns := []string{"id", "apelido", "nome", "nascimento", "stack", "search"}
	identifier := pgx.Identifier{"pessoas"}

	for {
		select {
		case payload := <-insertCh:
			_, err := w.db.CopyFrom(
				context.Background(),
				identifier,
				columns,
				pgx.CopyFromSlice(len(payload), w.makeCopyFromSlice(payload)),
			)

			if err != nil {
				log.Errorf("Error on insert batch: %v", err)
			}
		}
	}
}

func NewJobQueue() JobQueue {
	return make(JobQueue, MaxQueue)
}

func (Worker) makeCopyFromSlice(batch []Job) func(i int) ([]interface{}, error) {
	return func(i int) ([]interface{}, error) {
		return []interface{}{
			batch[i].Payload.ID,
			batch[i].Payload.Apelido,
			batch[i].Payload.Nome,
			batch[i].Payload.Nascimento,
			batch[i].Payload.StackStr(),
			batch[i].Payload.SearchStr(),
		}, nil
	}
}

func NewWorker(db *pgxpool.Pool, workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		db:         db,
	}
}

func NewDispatcher(db *pgxpool.Pool, jobQueue JobQueue) *Dispatcher {
	maxWorkers := MaxWorkers
	pool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		maxWorkers: maxWorkers,
		WorkerPool: pool,
		jobQueue:   jobQueue,
		db:         db,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.db, d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}
