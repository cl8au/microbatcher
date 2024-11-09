# Micro-batching library

This is technical interview task which implements a micro batching library. I just named it as **micro batcher**.

## Requirements

Micro-batching is a technique used in processing pipelines where individual tasks are grouped
together into small batches. This can improve throughput by reducing the number of requests made
to a downstream system. Your task is to implement a micro-batching library, with the following
requirements:

- it should allow the caller to submit a single Job, and it should return a JobResult
- it should process accepted Jobs in batches using a BatchProcessor
  - Don't implement BatchProcessor. This should be a dependency of your library.
- it should provide a way to configure the batching behaviour i.e. size and frequency
- it should expose a shutdown method which returns after all previously accepted Jobs are processed

## Design

- Each [batcher](https://github.com/cl8au/microbatcher/blob/main/batcher.go) is a worker which self contains the queue with batch size and timer in order to achieve micro batching processing.
- Giving library users flexibility to spawn multiple batchers if needed but also the control of job distribution.
- Job result and Job binding with field `ID` and generic in both Job and Job result should be able to handle multiple formats.
- Batch frequency and batch size are configurable and treated as inputs for batcher.
- Expose an additional method called `GetCurrentResults` which allows user to get current results at any points of time.
- `Shutdown` drains the job queue and I decide let user to call `GetCurrentResults` in order to get results rather than `Shutdown` returns the result in order to keep the syntax consistency with `Start` method.
- Contract of batch processor giving user flexibility to define what need to be return. Users can decide certain logic as dropping few job if needed.

## High level project structure

- pkg
  - contains all the business logic modules of zendesk search application
- vendor
  - contains all the saved dependencies resources which have been defined in go.mod
- go.mod
  - module definitions with all the dependencies
- batcher.go
  - core logic about batcher which does micro-batch processing
- Makefile
  - contains set of tasks which can build, test and lint this project
- `.golangci.yml`
  - contains all the lint configurations
- `.tools-versions`
  - `asdf` languages definition file

## Local setup

### Go version

Please install [asdf](https://asdf-vm.com/) and install defined Go version in `.tools-versions`. You can run `asdf install` if you have not installed this Go version yet,

### Build (build, lint and tests)

You can run `make` to run build which includes build, lint and tests. You also can see the test coverage in the stdout.

### Test coverage

You can run `make cov-html` to generate the html version of code coverage and `open cover.html` in order to browse the details. I have also committed the cover.html just to show the coverage. Please ignore `playground` folder since this is just for manual tests locally.

### Playground

You can find [playground](https://github.com/cl8au/microbatcher/blob/e468f6231807020f4d6aab7aed9a307886137ee7/playground/main.go) which does some manual tests apart from test coverage.
