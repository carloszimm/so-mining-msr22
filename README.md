# Stack Mining
Stack Overflow mining scripts used for the paper:
> Mining the Usage of Reactive Programming APIs: A Mining Study on GitHub and Stack Overflow.

## Data
Under the folders in `/assets`, data either genereated by or collected for the scripts execution can be found. The table gives a brief description of each folder:

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

Each result is comprised of three CSV files following the bellow file name pattern:
* [file name of the posts file]\_doctopicdist\_[#topics]\_[analyzed post field].csv
* [file name of the posts file]\_topicdist\_[#topics]\_[analyzed post field].csv
* [file name of the posts file]\_topicdist\_[#topics]\_[analyzed post field] - topwords.csv

Where:
* [file name of the posts file]: is a file under `assets/data explorer/consolidated sources` and set through [config](#configuration);
* [#topics]: number of topics for that specific execution;
* [analyzed post field]: either Title or Body (see [Configuration](#configuration)).

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

### Configuration
The LDA script require the setting of some configuration in a JSON(config.json) under `/configs` folder. This JSON is expecting a array of objects, each one representing a LDA execution. The objective must have the following structure (this is the object present by default in config.json):
```yaml
{
    "fileName": "all_withAnswers",
    "field": "Body",
    "combineTitleBody": true,
    "minTopics": 10,
    "maxTopics": 35,
    "sampleWords": 20
  }
```
Where:
* fileName(string): the name of the file with the posts(at `assets/data explorer/consolidated sources`);
* field(string): the field to considered in LDA (either Title or Body);
* combineTitleBody(boolean): set it to combine title and body and assign the result to the post's Body field (only applicable if `field` is set to `"Body"`);
* minTopics(integer): the minimum quantity of posts to be generated;
* maxTopics(integer): the maximum quantity of posts to be generated;
* sampleWords(integer): the amount of sample *top* words to be included in an extra file with file name ending with ` - topwords`.

#### Stack Exchange Data Explorer
Possible requirements:
* Internet browser
* Node.js (tested with v14.17.5)

We elaborated a tiny JS script to download the Stack Overflow posts (questions with and without accepted answers) related to the rx libraries from [Stack Exchange Data Explorer](https://data.stackexchange.com/) (SEDE).
It's available at `/scripts/data explorer/data-explorer.js`. To execute it, one must: (1) be logged in SEDE, (2) place the script in the DevTools's **Console**, and
(3) call `executeQuery` passing 0 (for RxJava), 1 (for RxJS), and 2 (for RxSwift) as a parameter. Moreover, there's a second script(`/scripts/data explorer/rename.js`) that can be used to move (and rename) the results to the their proper folder `/assets/data explorer/[rx library folder]`, so they can be further used by the `cmd/consolidate-sources/main.go` script. In order for this second JS script to work, one must place the results under `/scripts/data explorer/staging area` and call the script in a terminal (with node) and passing either 0 (for RxJava), 1 (for RxJS), and 2 (for RxSwift).
As detailed in the paper, these were the Stack Overflow tags used:
* *rx-java*, *rx-java2*, *rx-java3* (RxJava)
* *rxjs*, *rxjs5*, *rxjs6*, *rxjs7* (RxJS)
* *rx-swift* (RxSwift)
