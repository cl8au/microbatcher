package types

type JobId interface {
	~int | ~string
}

type Job[I JobId, T any] struct {
	ID   I
	Data T
}

type JobResult[I JobId, T any] struct {
	ID     I
	Result T
	Errors error
}
