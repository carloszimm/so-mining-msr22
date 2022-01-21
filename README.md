# Stack Mining
Stack Overflow mining scripts used for the paper:
> Mining the Usage of Reactive Programming API: A Mining Study on GitHub and Stack Overflow.

## Requirements
Most of the scripts utilize Golang as the main language and they have be executed the following version:
* Go version 1.17.5

## Execution
The Go scripts are available under the `/cmd` folder
```go
go run cmd/consolidate-sources/main.go
````
:computer: Script to unify all the CSV acquired from [Stack Exchange Data Explorer](https://data.stackexchange.com/).

:floppy_disk: After execution, the result is available at `assets/data explorer/consolidated sources/`.
```go
go run cmd/extract-posts/main.go
```
:computer: Script to extract post from a given topic.

:floppy_disk: After execution, the result is available at `assets/extracted-posts`.
```go
go run cmd/lda/main.go
```
:computer: Script to execute the LDA algorithm.

:floppy_disk: After execution, the result is available at `assets/lda-results`.
```go
go run cmd/open-sort/main.go
```
:computer: Script to generate random posts according to their topics and facilitate the open sort (topic labeling) execution.

:floppy_disk: After execution, the result is available at `assets/opensort`.
```go
go run cmd/operators-search/main.go
```
:computer: Script to search for operators among the Stack Overflow posts.

:floppy_disk: After execution, the result is available at `assets/operators-search`.
```go
go run cmd/process-results/main.go
```
:computer: Script to process results and generate info about the topics, the popularities and difficulties.

:floppy_disk: After execution, the result is available at `assets/result-processing`.
