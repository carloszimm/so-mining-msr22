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
* **_[file name of the posts file]_\_doctopicdist\__[#topics]_\__[analyzed post field]_.csv** - contains the posts' ids and their distribution of topics+proportion, including the dominant topic and its proportion in a separate column for easy retrieval;
* **_[file name of the posts file]_\_topicdist\__[#topics]_\__[analyzed post field]_.csv** - the topic distribution along with their words+proportion descendingly sorted by word proportion;
* **_[file name of the posts file]_\_topicdist\__[#topics]_\__[analyzed post field]_ - topwords.csv** - (extra) the same as the above one but presenting the topics only with their top words (set in [config](#configuration)) to facilitate the open card sorting technique.

Where:
* _[file name of the posts file]_: is a file under `assets/data explorer/consolidated sources` and set through [config](#configuration);
* _[#topics]_: number of topics for that specific execution;
* _[analyzed post field]_: either Title or Body (see [Configuration](#configuration)).

## Execution
### Requirements
Most of the scripts utilize Golang as the main language and they have be executed the following version:
* Go v1.17.5

Before execution of the Golang scripts, the following command must be issued in a terminal (inside the root of the project) to download the dependencies:
```sh
go mod tidy
```

### Scripts
The Go scripts are available under the `/cmd` folder

##### consolidate-sources
Script to unify all the CSV acquired from [Stack Exchange Data Explorer](https://data.stackexchange.com/).
```sh
go run cmd/consolidate-sources/main.go
````
&ensp;:floppy_disk: After execution, the result is available at `assets/data explorer/consolidated sources/`.

##### extract-posts
Script to extract post from a given topic.
```sh
go run cmd/extract-posts/main.go
```
&ensp;:floppy_disk: After execution, the result is available at `assets/extracted-posts`.

##### lda
Script to execute the LDA algorithm.
```sh
go run cmd/lda/main.go
``` 
&ensp;:floppy_disk: After execution, the result is available at `assets/lda-results`.

##### open-sort
Script to generate random posts according to their topics and facilitate the open sort (topic labeling) execution.
```sh
go run cmd/open-sort/main.go
```
&ensp;:floppy_disk: After execution, the result is available at `assets/opensort`.

##### operators-search
Script to search for operators among the Stack Overflow posts.
```sh
go run cmd/operators-search/main.go
```
&ensp;:floppy_disk: After execution, the result is available at `assets/operators-search`.

##### process-results
Script to process results and generate info about the topics, the popularities and difficulties.
```sh
go run cmd/process-results/main.go
```
&ensp;:floppy_disk: After execution, the result is available at `assets/result-processing`.

#### Configuration
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
* **fileName(string)**: the name of the file with the posts(at `assets/data explorer/consolidated sources`);
* **field(string)**: the field to considered in LDA (either Title or Body);
* **combineTitleBody(boolean)**: set it to combine title and body and assign the result to the post's Body field (only applicable if `field` is set to `"Body"`);
* **minTopics(integer)**: the minimum quantity of posts to be generated;
* **maxTopics(integer)**: the maximum quantity of posts to be generated;
* **sampleWords(integer)**: the amount of sample *top* words to be included in an extra file with file name ending with ` - topwords`.

#### Stack Exchange Data Explorer
Possible requirements:
* Internet browser
* Node.js (tested with v14.17.5)

We elaborated a tiny JS script to download the Stack Overflow posts (questions with and without accepted answers) related to the rx libraries from [Stack Exchange Data Explorer](https://data.stackexchange.com/) (SEDE).
It's available at `/scripts/data explorer/data-explorer.js`. To execute it, one must:
1. Be logged in SEDE;
2. Place the script in the DevTools's **Console**;
3. Call `executeQuery` passing 0 (for RxJava), 1 (for RxJS), and 2 (for RxSwift) as a parameter.

Moreover, there's a second script(`/scripts/data explorer/rename.js`) that can be used to move (and rename) the results to the their proper folder `/assets/data explorer/[rx library folder]`, so they can be further used by the `cmd/consolidate-sources/main.go` script. In order for this second JS script to work, one must place the results under `/scripts/data explorer/staging area` and call the script in a terminal (with node) and passing either 0 (for RxJava), 1 (for RxJS), and 2 (for RxSwift). For example:
```sh
node rename 0
```
Before execution of node.js script, one must execute the following terminal command within `/scripts/data explorer/`:
```sh
npm install
```
As detailed in the paper, these were the Stack Overflow tags used:
* *rx-java*, *rx-java2*, *rx-java3* (RxJava)
* *rxjs*, *rxjs5*, *rxjs6*, *rxjs7* (RxJS)
* *rx-swift* (RxSwift)

## Other Useful Information
### Stack Overflow Removed Terms
As defined in the preprocessing phase in the paper, some terms commonly found in the Stack Overflow posts were removed from the corpus. Those include:
<blockquote>
<code>differ</code>, <code>specif</code>, <code>deal</code>, <code>prefer</code>, <code>easili</code>, <code>easier</code>,
<code>mind</code>, <code>current</code>, <code>solv</code>, <code>proper</code>, <code>modifi</code>, <code>explain</code>,
<code>hope</code>, <code>help</code>, <code>wonder</code>, <code>altern</code>, <code>sens</code>, <code>entir</code>,
<code>ps</code>, <code>solut</code>, <code>achiev</code>, <code>approach</code>, <code>answer</code>, <code>requir</code>,
<code>lot</code>, <code>feel</code>, <code>pretti</code>, <code>easi</code>, <code>goal</code>, <code>think</code>,
<code>complex</code>, <code>eleg</code>, <code>improv</code>, <code>look</code>, <code>complic</code>, <code>day</code>,
<code>chang</code>, <code>issu</code>, <code>add</code>, <code>edit</code>, <code>remov</code>, <code>custom</code>,
<code>suggest</code>, <code>comment</code>, <code>ad</code>, <code>refer</code>, <code>stackblitz</code>, <code>link</code>,
<code>mention</code>, <code>detect</code>, <code>face</code>, <code>fix</code>, <code>attach</code>, <code>perfect</code>,
<code>mark</code>, <code>reason</code>, <code>suppos</code>, <code>notic</code>, <code>snippet</code>, <code>demo</code>,
<code>line</code>, <code>piec</code>, <code>appear</code>
</blockquote>

### Topic-Label Mapping
| Topic #      | Label/Name    |
| ------------ |:-------------|
| 0 | Concurrency |
| 1 | Stream Creation and Composition |
| 2 | Typing and Correctness |
| 3 | UI for Web-based Systems |
| 4 | Input Validation |
| 5 | Introductory Questions |
| 6 | Testing and Debugging |
| 7 | REST API Calls |
| 8 | Android Development |
| 9 | Data Access |
| 10 | State Management and JavaScript |
| 11 | Control Flow |
| 12 | HTTP Handling |
| 13 | Stream Manipulation |
| 14 | Error Handling |
| 15 | Stream Lifecycle |
| 16 | Array Manipulation |
| 17 | Web Development |
| 18 | General Programming |
| 19 | iOS Development |
| 20 | Multicasting |
| 21 | Timing |
| 22 | Dependency Management  |
