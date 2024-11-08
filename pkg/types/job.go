package types

import "fmt"

type JobId interface {
	~int | ~string
}

type Job[I JobId, T any] struct {
	ID   I
	Data T
}

func (j *Job[I, T]) String() string {
	return fmt.Sprintf("job: id=%v", j.ID)
}

type JobResult[I JobId, T any] struct {
	ID     I
	Data   T
	Errors error
}

func (jr *JobResult[I, T]) String() string {
	return fmt.Sprintf("job result: id=%v", jr.ID)
}
