[![Go Report Card](https://goreportcard.com/badge/github.com/masonj188/ganalyze)](https://goreportcard.com/report/github.com/masonj188/ganalyze)

# Ganalyze
Ganalyze generates a set of html files that give information about given PE files and use an ML model to make a determination on whether or not the file is malicious or benign.

## Building from Source

In order to clone this repo, you must have [git-lfs](https://git-lfs.github.com/) installed.

Building and running Ganalyze should be done in a Python 3.6 virtual environment. Run `go build` from the repository's root directory, then `cd python` and run `./installember.sh`.

## Usage

To use the program, run `./ganalyze <directory>`, with directory being the root directory of the PE files you'd like to analyze.  Ganalyze will recursively scan through the directory and analyze each PE file. A `report` directory is produced which contains an index.html file and html files for each PE file analyzed. Index.html contains a listing of all the files analyzed and links to each of the respective html reports.  The reports directory will maintain the directory structure of the original directory passed in to Ganalyze.
