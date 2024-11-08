package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobStringWithIntegerId(t *testing.T) {
	tests := []struct {
		name     string
		job      *Job[int, string]
		expected string
	}{
		{
			name: "Test with integer Id",
			job: &Job[int, string]{
				ID:   123,
				Data: "data 123 with int id",
			},
			expected: "job: id=123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.job.String())
		})
	}
}

func TestJobStringWithStringId(t *testing.T) {
	tests := []struct {
		name     string
		job      *Job[string, string]
		expected string
	}{
		{
			name: "Test with integer String",
			job: &Job[string, string]{
				ID:   "123",
				Data: "data 123 with string id",
			},
			expected: "job: id=123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.job.String())
		})
	}
}

func TestJobResultStringWithIntegerId(t *testing.T) {
	tests := []struct {
		name      string
		jobResult *JobResult[int, string]
		expected  string
	}{
		{
			name: "Test with integer Id",
			jobResult: &JobResult[int, string]{
				ID:     123,
				Data:   "result data with int id",
				Errors: nil,
			},
			expected: "job result: id=123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.jobResult.String())
		})
	}
}

func TestJobResultStringWithStringId(t *testing.T) {
	tests := []struct {
		name      string
		jobResult *JobResult[string, string]
		expected  string
	}{
		{
			name: "Test with string Id",
			jobResult: &JobResult[string, string]{
				ID:     "123",
				Data:   "result data with string id",
				Errors: nil,
			},
			expected: "job result: id=123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.jobResult.String())
		})
	}
}
