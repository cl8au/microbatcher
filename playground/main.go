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

	// first round
	_ = batcher01.Start()
	_ = batcher02.Start()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_, _ = batcher01.Submit(&types.Job[int, string]{
				ID:   i,
				Data: fmt.Sprintf("first round job %d", i),
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			_, _ = batcher02.Submit(&types.Job[int, string]{
				ID:   i,
				Data: fmt.Sprintf("first round job %d", i),
			})
		}
	}()

	wg.Wait()

	_ = batcher01.Shutdown()
	_ = batcher02.Shutdown()

	results := batcher01.GetCurrentResults()
	for _, result := range results {
		slog.Info(fmt.Sprintf("first round batch01 result: %s with data: %s\n", result, result.Data))
	}
	slog.Info(fmt.Sprintf("first round batch01 result size: %d\n", len(results)))
	results = batcher02.GetCurrentResults()
	for _, result := range results {
		slog.Info(fmt.Sprintf("first round batch02 result: %s with data: %s\n", result, result.Data))
	}
	slog.Info(fmt.Sprintf("first round batch02 result size: %d\n", len(results)))

	// second round
	_ = batcher01.Start()
	_ = batcher02.Start()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_, _ = batcher01.Submit(&types.Job[int, string]{
				ID:   i,
				Data: fmt.Sprintf("second round job %d", i),
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			_, _ = batcher02.Submit(&types.Job[int, string]{
				ID:   i,
				Data: fmt.Sprintf("second round job %d", i),
			})
		}
	}()

	wg.Wait()

	_ = batcher01.Shutdown()
	_ = batcher02.Shutdown()

	results = batcher01.GetCurrentResults()
	for _, result := range results {
		slog.Info(fmt.Sprintf("second round batch01 result: %s with data: %s\n", result, result.Data))
	}
	slog.Info(fmt.Sprintf("second round batch01 result size: %d\n", len(results)))
	results = batcher02.GetCurrentResults()
	for _, result := range results {
		slog.Info(fmt.Sprintf("second round batch02 result: %s with data: %s\n", result, result.Data))
	}
	slog.Info(fmt.Sprintf("second round batch02 result size: %d\n", len(results)))
}
