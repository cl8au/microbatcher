package configs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	assert.Equal(t, DEFAULT_QUEUE_SIZE, config.GetJobQueueSize())
	assert.Equal(t, DEFAULT_BATCH_PROCESS_SIZE, config.GetBatchProcessSize())
	assert.Equal(t, DEFAULT_BATCH_PROCESS_FREQUENCY_IN_MILLISECOND*time.Millisecond, config.GetBatchProcessFrequency())
}

func TestNewCustomConfig(t *testing.T) {
	tests := []struct {
		name                string
		jobQueueSize        int
		batchProcessSize    int
		batchProcessFreq    time.Duration
		expectError         bool
		expectedErrorString string
	}{
		{
			name:             "Valid custom batcher config",
			jobQueueSize:     500,
			batchProcessSize: 200,
			batchProcessFreq: 200 * time.Second,
			expectError:      false,
		},
		{
			name:                "Invalid custom batcher config by zero job queue size",
			jobQueueSize:        0,
			batchProcessSize:    200,
			batchProcessFreq:    200 * time.Second,
			expectError:         true,
			expectedErrorString: "jobQueueSize, batchProcessSize, and batchProcessFrequency must be positive",
		},
		{
			name:                "Invalid custom batcher config by zero batch process frequency",
			jobQueueSize:        500,
			batchProcessSize:    100,
			batchProcessFreq:    0,
			expectError:         true,
			expectedErrorString: "jobQueueSize, batchProcessSize, and batchProcessFrequency must be positive",
		},
		{
			name:                "Invalid custom batcher config by zero batch process size",
			jobQueueSize:        500,
			batchProcessSize:    0,
			batchProcessFreq:    200 * time.Second,
			expectError:         true,
			expectedErrorString: "jobQueueSize, batchProcessSize, and batchProcessFrequency must be positive",
		},
		{
			name:                "Invalid custom batcher config by job queue size too small",
			jobQueueSize:        200,
			batchProcessSize:    200,
			batchProcessFreq:    200 * time.Second,
			expectError:         true,
			expectedErrorString: "job queue size must be at least twice the batch process size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewCustomConfig(tt.jobQueueSize, tt.batchProcessSize, tt.batchProcessFreq)

			if tt.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErrorString)
				assert.Equal(t, BatcherConfig{}, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tt.jobQueueSize, config.GetJobQueueSize())
				assert.Equal(t, tt.batchProcessSize, config.GetBatchProcessSize())
				assert.Equal(t, tt.batchProcessFreq, config.GetBatchProcessFrequency())
			}
		})
	}
}
