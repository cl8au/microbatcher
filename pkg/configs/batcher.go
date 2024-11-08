package configs

import (
	"errors"
	"time"
)

const DEFAULT_QUEUE_SIZE = 100
const DEFAULT_BATCH_PROCESS_SIZE = 10
const DEFAULT_BATCH_PROCESS_FREQUENCY_IN_MILLISECOND = 100
const QUEUE_FACTOR = 2

type BatcherConfig struct {
	jobQueueSize          int
	batchProcessSize      int
	batchProcessFrequency time.Duration
}

// NewDefaultConfig creates and returns a new batcher config with default values.
func NewDefaultConfig() BatcherConfig {
	return BatcherConfig{
		jobQueueSize:          DEFAULT_QUEUE_SIZE,
		batchProcessSize:      DEFAULT_BATCH_PROCESS_SIZE,
		batchProcessFrequency: DEFAULT_BATCH_PROCESS_FREQUENCY_IN_MILLISECOND * time.Millisecond,
	}
}

// NewCustomConfig creates and returns a new batcher with custom values.
func NewCustomConfig(
	jobQueueSize int,
	batchProcessSize int,
	batchProcessFrequency time.Duration,
) (BatcherConfig, error) {
	if jobQueueSize < 1 || batchProcessSize < 1 || batchProcessFrequency <= 0 {
		return BatcherConfig{}, errors.New("jobQueueSize, batchProcessSize, and batchProcessFrequency must be positive")
	}

	if jobQueueSize < batchProcessSize*QUEUE_FACTOR {
		return BatcherConfig{}, errors.New("job queue size must be at least twice the batch process size")
	}

	return BatcherConfig{
		jobQueueSize:          jobQueueSize,
		batchProcessSize:      batchProcessSize,
		batchProcessFrequency: batchProcessFrequency,
	}, nil
}

// GetJobQueueSize returns the job queue size.
func (b *BatcherConfig) GetJobQueueSize() int {
	return b.jobQueueSize
}

// GetBatchProcessSize returns the batch process size.
func (b *BatcherConfig) GetBatchProcessSize() int {
	return b.batchProcessSize
}

// GetBatchProcessFrequency returns the batch process frequency.
func (b *BatcherConfig) GetBatchProcessFrequency() time.Duration {
	return b.batchProcessFrequency
}
