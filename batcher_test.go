package microbatcher

import (
	"fmt"
	"microbatcher/pkg/configs"
	"microbatcher/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestingMicroBatcherProcess[I types.JobId] struct {
	slow          bool
	sleepDuration time.Duration
}

func (tm *TestingMicroBatcherProcess[I]) Process(jobs []*types.Job[I, string]) []*types.JobResult[I, string] {
	results := make([]*types.JobResult[I, string], 0)

	if tm.slow {
		time.Sleep(tm.sleepDuration)
	}

	for _, job := range jobs {
		results = append(results, &types.JobResult[I, string]{
			ID:   job.ID,
			Data: fmt.Sprintf("%v is processed", job.ID),
		})
	}
	return results
}

var jobs = []*types.Job[string, string]{
	{
		ID:   "job1",
		Data: "data1",
	},
	{
		ID:   "job2",
		Data: "data2",
	},
	{
		ID:   "job3",
		Data: "data3",
	},
}

var expectedJobResults = map[string]*types.JobResult[string, string]{
	"job1": {
		ID:   "job1",
		Data: "job1 is processed",
	},
	"job2": {
		ID:   "job2",
		Data: "job2 is processed",
	},
	"job3": {
		ID:   "job3",
		Data: "job3 is processed",
	},
}

func TestMicroBatcherProcessBySize(t *testing.T) {
	tests := []struct {
		name                string
		config              configs.BatcherConfig
		jobs                []*types.Job[string, string]
		expectedJobResults  map[string]*types.JobResult[string, string]
		expectError         bool
		expectedErrorString string
	}{
		{
			name: "Submits jobs with exceeds the batch process size",
			config: func() configs.BatcherConfig {
				cfg, _ := configs.NewCustomConfig(10, 3, 5*time.Second)
				return cfg
			}(),
			jobs:               jobs,
			expectedJobResults: expectedJobResults,
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := NewMicroBatcher("tester", &TestingMicroBatcherProcess[string]{}, tt.config)
			startErr := mb.Start()
			assert.Nil(t, startErr)

			for _, job := range tt.jobs {
				result, err := mb.Submit(job)
				if tt.expectError {
					assert.Error(t, err)
					assert.EqualError(t, err, tt.expectedErrorString)
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Nil(t, err)
					assert.Equal(t, job.ID, result.ID)
				}
			}

			time.Sleep(2 * time.Second)

			shutdownErr := mb.Shutdown()
			assert.Nil(t, shutdownErr)
			results := mb.GetCurrentResults()
			assert.Equal(t, len(results), len(tt.expectedJobResults))

			for _, result := range results {
				expectedResult, ok := tt.expectedJobResults[result.ID]
				assert.True(t, ok)
				assert.Equal(t, expectedResult.Data, result.Data)
			}
		})
	}
}

func TestMicroBatcherProcessByTimer(t *testing.T) {
	tests := []struct {
		name                string
		config              configs.BatcherConfig
		jobs                []*types.Job[string, string]
		expectedJobResults  map[string]*types.JobResult[string, string]
		expectError         bool
		expectedErrorString string
	}{
		{
			name: "Submits jobs with which is processed by timer",
			config: func() configs.BatcherConfig {
				cfg, _ := configs.NewCustomConfig(100, 20, 1*time.Second)
				return cfg
			}(),
			jobs:               jobs,
			expectedJobResults: expectedJobResults,
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := NewMicroBatcher("tester", &TestingMicroBatcherProcess[string]{}, tt.config)
			startErr := mb.Start()
			assert.Nil(t, startErr)

			for _, job := range tt.jobs {
				result, err := mb.Submit(job)
				if tt.expectError {
					assert.Error(t, err)
					assert.EqualError(t, err, tt.expectedErrorString)
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Nil(t, err)
					assert.Equal(t, job.ID, result.ID)
				}
			}

			time.Sleep(2 * time.Second)

			shutdownErr := mb.Shutdown()
			assert.Nil(t, shutdownErr)
			results := mb.GetCurrentResults()
			assert.Equal(t, len(results), len(tt.expectedJobResults))

			for _, result := range results {
				expectedResult, ok := tt.expectedJobResults[result.ID]
				assert.True(t, ok)
				assert.Equal(t, expectedResult.Data, result.Data)
			}
		})
	}
}

func TestMicroBatcherProcessBySizeStartAndShutdownMultipleTimes(t *testing.T) {
	tests := []struct {
		name                string
		config              configs.BatcherConfig
		jobs                []*types.Job[string, string]
		expectedJobResults  map[string]*types.JobResult[string, string]
		expectError         bool
		expectedErrorString string
	}{
		{
			name: "Submits jobs with exceeds the batch process size with multiple times",
			config: func() configs.BatcherConfig {
				cfg, _ := configs.NewCustomConfig(10, 3, 5*time.Second)
				return cfg
			}(),
			jobs:               jobs,
			expectedJobResults: expectedJobResults,
			expectError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := NewMicroBatcher("tester", &TestingMicroBatcherProcess[string]{}, tt.config)
			for i := 0; i < 2; i++ {
				startErr := mb.Start()
				assert.Nil(t, startErr)
				assert.Empty(t, mb.GetCurrentResults())
				for _, job := range tt.jobs {
					result, err := mb.Submit(job)
					if tt.expectError {
						assert.Error(t, err)
						assert.EqualError(t, err, tt.expectedErrorString)
						assert.Nil(t, result)
					} else {
						assert.NotNil(t, result)
						assert.Nil(t, err)
						assert.Equal(t, job.ID, result.ID)
					}
				}

				time.Sleep(2 * time.Second)

				shutdownErr := mb.Shutdown()
				assert.Nil(t, shutdownErr)
				results := mb.GetCurrentResults()
				assert.Equal(t, len(results), len(tt.expectedJobResults))

				for _, result := range results {
					expectedResult, ok := tt.expectedJobResults[result.ID]
					assert.True(t, ok)
					assert.Equal(t, expectedResult.Data, result.Data)
				}
			}
		})
	}
}

func TestMicroBatcherStartError(t *testing.T) {
	tests := []struct {
		name   string
		config configs.BatcherConfig
	}{
		{
			name:   "Start batcher multiple times",
			config: configs.NewDefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := NewMicroBatcher("tester", &TestingMicroBatcherProcess[string]{}, tt.config)
			startErr := mb.Start()
			assert.Nil(t, startErr)

			// start again to get error since it is running already
			startErr = mb.Start()
			assert.EqualError(t, startErr, "batcher is started already")
		})
	}
}

func TestMicroBatcherShutdownError(t *testing.T) {
	tests := []struct {
		name   string
		config configs.BatcherConfig
	}{
		{
			name:   "Shutdown batcher which is stopped already",
			config: configs.NewDefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := NewMicroBatcher("tester", &TestingMicroBatcherProcess[string]{}, tt.config)
			shutdownErr := mb.Shutdown()
			assert.NotNil(t, shutdownErr)
			assert.EqualError(t, shutdownErr, "invalid shutdown since batcher is stopped")
		})
	}
}

func TestMicroBatcherSubmitErrorSinceNotStart(t *testing.T) {
	tests := []struct {
		name   string
		config configs.BatcherConfig
	}{
		{
			name:   "Submits a job to stopped batcher",
			config: configs.NewDefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := NewMicroBatcher("tester", &TestingMicroBatcherProcess[string]{}, tt.config)
			result, submitErr := mb.Submit(&types.Job[string, string]{})
			assert.Nil(t, result)
			assert.NotNil(t, submitErr)
			assert.EqualError(t, submitErr, "invalid submission since batcher is not started")
		})
	}
}

func TestMicroBatcherSubmitErrorSinceQueueIsFull(t *testing.T) {
	tests := []struct {
		name   string
		config configs.BatcherConfig
	}{
		{
			name: "Submits a job when job queue is full",
			config: func() configs.BatcherConfig {
				cfg, _ := configs.NewCustomConfig(2, 1, 100*time.Second)
				return cfg
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := NewMicroBatcher("tester", &TestingMicroBatcherProcess[string]{
				slow:          true,
				sleepDuration: 1 * time.Minute,
			}, tt.config)

			_ = mb.Start()

			_, submitOneErr := mb.Submit(&types.Job[string, string]{
				ID:   "job11",
				Data: "job11 is processed",
			})
			assert.Nil(t, submitOneErr)
			_, submitTwoErr := mb.Submit(&types.Job[string, string]{
				ID:   "job12",
				Data: "job12 is processed",
			})
			assert.Nil(t, submitTwoErr)
			_, submitThreeErr := mb.Submit(&types.Job[string, string]{
				ID:   "job13",
				Data: "job13 is processed",
			})
			assert.NotNil(t, submitThreeErr)
			assert.EqualError(t, submitThreeErr, "job queue is full")
		})
	}
}
