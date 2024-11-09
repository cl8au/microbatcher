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

## Local setup

### Go version

Please install [asdf](https://asdf-vm.com/) and install defined Go version in `.tools-versions`. You can run `asdf install` if you have not installed this Go version yet,

### Build (build, lint and tests)

You can run `make` to run build which includes build, lint and tests. You also can see the test coverage in the stdout.

### Test coverage

You can run `make cov-html` to generate the html version of code coverage and `open cover.html` in order to browse the details. I have also committed the cover.html just to show the coverage. Please ignore `playground` folder since this is just for manual tests locally.
