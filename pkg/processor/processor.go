package processor

import "microbatcher/pkg/types"

// BatchProcessor is a contract for processing batch of jobs and return results.
// You will need to implement your specific business process logic.
type BatchProcessor[I types.JobId, T any] interface {
	Process(jobs []*types.Job[I, T]) []*types.JobResult[I, T]
}
