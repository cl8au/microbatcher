package microbatcher

import (
	"errors"
	"fmt"
	"log/slog"
	"microbatcher/pkg/configs"
	"microbatcher/pkg/processor"
	"microbatcher/pkg/types"
	"sync"
	"time"
)

type microBatcher[I types.JobId, T any] struct {
	name         string
	processor    processor.BatchProcessor[I, T]
	config       configs.BatcherConfig
	results      []*types.JobResult[I, T]
	resultsMutex sync.Mutex
	running      bool
	runningMutex sync.Mutex
	jobs         chan *types.Job[I, T]
	shutdown     chan struct{}
	wg           sync.WaitGroup
}

// NewMicroBatcher creates a new instance of the micro batcher with processor and configurations.
// Each batcher is like a batch worker, it retain it's own job queue and results.
func NewMicroBatcher[I types.JobId, T any](
	name string,
	processor processor.BatchProcessor[I, T],
	config configs.BatcherConfig,
) *microBatcher[I, T] {
	return &microBatcher[I, T]{
		name:      name,
		processor: processor,
		config:    config,
	}
}

// Submit submits a new job to the internal job queue and returns a job result with accepted state.
func (mb *microBatcher[I, T]) Submit(job *types.Job[I, T]) (*types.JobResult[I, T], error) {
	mb.runningMutex.Lock()
	defer mb.runningMutex.Unlock()

	if !mb.running {
		return nil, errors.New("invalid submission since batcher is not started")
	}

	select {
	case mb.jobs <- job:
		slog.Info(fmt.Sprintf("%s submits %s", mb.name, job))
		return &types.JobResult[I, T]{ID: job.ID}, nil
	default:
		return nil, errors.New("job queue is full")
	}
}

// Start starts the batch process goroutine which execute custom processor either by either timer
// or size constraint
func (mb *microBatcher[I, T]) Start() error {
	mb.runningMutex.Lock()
	defer mb.runningMutex.Unlock()

	if mb.running {
		return errors.New("batcher is started already")
	}

	slog.Info(fmt.Sprintf("%s starts", mb.name))

	mb.running = true
	// init channels here. This is helpful to
	// let batcher can be shutdown and start again
	mb.jobs = make(chan *types.Job[I, T], mb.config.GetJobQueueSize())
	mb.shutdown = make(chan struct{})
	mb.results = nil

	mb.wg.Add(1)
	go mb.execute()
	return nil
}

// GetCurrentResults returns the all current of processed jobs
func (mb *microBatcher[I, T]) GetCurrentResults() []*types.JobResult[I, T] {
	mb.resultsMutex.Lock()
	defer mb.resultsMutex.Unlock()

	clonedResults := make([]*types.JobResult[I, T], len(mb.results))
	copy(clonedResults, mb.results)
	return clonedResults
}

func (mb *microBatcher[I, T]) Shutdown() error {
	mb.runningMutex.Lock()
	defer mb.runningMutex.Unlock()

	if !mb.running {
		return errors.New("invalid shutdown since batcher is stopped")
	}
	slog.Info(fmt.Sprintf("%s starts shutting down", mb.name))
	// send shutdown signal via channel
	close(mb.shutdown)
	// wait for the process goroutine to be finished
	mb.wg.Wait()
	// close job channels
	close(mb.jobs)

	mb.running = false
	slog.Info(fmt.Sprintf("%s shuts down", mb.name))

	return nil
}

func (mb *microBatcher[I, T]) execute() {
	defer mb.wg.Done()

	timer := time.NewTimer(mb.config.GetBatchProcessFrequency())
	defer timer.Stop()

	// rely on local batch job slice to monitor the in-taking batch size
	var batchJobs []*types.Job[I, T]
	for {
		select {
		case job := <-mb.jobs:
			batchJobs = append(batchJobs, job)
			// invoke custom processor when batch size is reached
			if len(batchJobs) >= mb.config.GetBatchProcessSize() {
				mb.processBatch(batchJobs, timer)
				batchJobs = nil
			}
		case <-timer.C:
			if len(batchJobs) > 0 {
				// Process batch on timer trigger
				mb.processBatch(batchJobs, timer)
				batchJobs = nil
			}
		case <-mb.shutdown:
			// handle shutdown case
			batchJobs = mb.drainQueue(batchJobs)
			if len(batchJobs) > 0 {
				mb.processBatch(batchJobs, timer)
			}
			return
		}
	}
}

func (mb *microBatcher[I, T]) processBatch(batchJobs []*types.Job[I, T], timer *time.Timer) {
	// need to stop and reset the timer since the batch process
	timer.Stop()
	defer timer.Reset(mb.config.GetBatchProcessFrequency())

	// call custom processor to process the batch jobs
	slog.Info(fmt.Sprintf("%s starts batch process", mb.name))
	results := mb.processor.Process(batchJobs)

	// cache this batch results in the batcher
	mb.recordResults(results)
}

// The mutex here since the GetCurrentResults function. Read and write in different goroutine and GetCurrentResults
// can be called by external anytime they need. Therefore, the mutex of results is needed here.
func (mb *microBatcher[I, T]) recordResults(newResults []*types.JobResult[I, T]) {
	mb.resultsMutex.Lock()
	defer mb.resultsMutex.Unlock()

	mb.results = append(mb.results, newResults...)
}

func (mb *microBatcher[I, T]) drainQueue(batchJobs []*types.Job[I, T]) []*types.Job[I, T] {
	for {
		select {
		case job := <-mb.jobs:
			batchJobs = append(batchJobs, job)
		default:
			// no more to drain
			return batchJobs
		}
	}
}
