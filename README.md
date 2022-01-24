# Stack Mining
Stack Overflow mining scripts used for the paper:
> Mining the Usage of Reactive Programming APIs: A Mining Study on GitHub and Stack Overflow.

## Data
Under the `/assets` folders, data either genereated by or collected for the scripts execution can be found. The table gives a brief description of each folder:

| Folder   | Description         |
| :------------- |:-------------|
| data explorer | Contains posts collected from Stack Exchange Data Explorer |
| extracted-posts | Includes JSON files having the posts related to the most relenvat topics (RQ3) |
| lda-results | Contains the results of the last LDA execution |
| operators-search | Includes the results for the operator search for Rx libraries |
| operators | Includes JSON files consisting of Rx libraries' operators |
| result-processing | Contains data presented in the Result section (RQ2) |

The file `stopwords.txt` contains a list of stop words used during preprocessing.

### LDA results
The results for the last LDA (Latent Dirichlet Allocation) are available under `/assets/2022-01-12 02-21-28/`. As detailed in the paper, the execution with the following settings generated the most coherent results:
| Parameter     | Value         |
| :------------- |:-------------:|
| Topic         | 23 |
| HyperParameters | &alpha;=&eta;=0.01 |
| Iterations | 1,000 |

## Execution
### Requirements
Most of the scripts utilize Golang as the main language and they have be executed the following version:
* Go version 1.17.5

### Scripts
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
