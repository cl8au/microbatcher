package main

import (
	"fmt"
	"log/slog"
	"microbatcher"
	"microbatcher/pkg/configs"
	"microbatcher/pkg/types"
	"sync"
)

type PlaygroundMicroBatcherProcess struct{}

func (pm *PlaygroundMicroBatcherProcess) Process(jobs []*types.Job[int, string]) []*types.JobResult[int, string] {
	results := make([]*types.JobResult[int, string], 0)
	for _, job := range jobs {
		slog.Info(fmt.Sprintf("playground is processing job %d", job.ID))
		results = append(results, &types.JobResult[int, string]{
			ID:   job.ID,
			Data: fmt.Sprintf("processed job %d", job.ID),
		})
	}
	return results
}

func main() {
	// create two batchers
	processor := &PlaygroundMicroBatcherProcess{}
	batcher01 := microbatcher.NewMicroBatcher("batcher01", processor, configs.NewDefaultConfig())
	batcher02 := microbatcher.NewMicroBatcher("batcher02", processor, configs.NewDefaultConfig())

	for k := 0; k < 1; k++ {
		startBatcher01Err := batcher01.Start()
		if startBatcher01Err != nil {
			slog.Error("failed to start batcher01")
		}
		startBatcher02Err := batcher02.Start()
		if startBatcher02Err != nil {
			slog.Error("failed to start batcher02")
		}

		var wg sync.WaitGroup

		// NOTE: expect some failures of submit since the queue is too small and job size is too large
		// This is desired just for testing purpose.
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				_, submitErr := batcher01.Submit(&types.Job[int, string]{
					ID:   i,
					Data: fmt.Sprintf("round %d - job %d", k, i),
				})
				if submitErr != nil {
					slog.Error(fmt.Sprintf("round %d - failed to submit job %d", k, i))
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				_, submitErr := batcher02.Submit(&types.Job[int, string]{
					ID:   i,
					Data: fmt.Sprintf("round %d - job %d", k, i),
				})
				if submitErr != nil {
					slog.Error(fmt.Sprintf("round %d - failed to submit job %d", k, i))
				}
			}
		}()

		wg.Wait()

		shutdownBatcher01Err := batcher01.Shutdown()
		if shutdownBatcher01Err != nil {
			slog.Error("failed to start batcher01")
		}
		shutdownBatcher02Err := batcher02.Shutdown()
		if shutdownBatcher02Err != nil {
			slog.Error("failed to start batcher02")
		}

		results := batcher01.GetCurrentResults()
		for _, result := range results {
			slog.Info(fmt.Sprintf("round %d - batch01 result: %s with data: %s\n", k, result, result.Data))
		}
		slog.Info(fmt.Sprintf("round %d - batch01 result size: %d\n", k, len(results)))
		results = batcher02.GetCurrentResults()
		for _, result := range results {
			slog.Info(fmt.Sprintf("round %d - batch02 result: %s with data: %s\n", k, result, result.Data))
		}
		slog.Info(fmt.Sprintf("round %d - batch02 result size: %d\n", k, len(results)))
	}
}
