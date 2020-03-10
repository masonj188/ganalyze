# Ganalyze
Ganalyze generates a set of html files that give information about given PE files and use an ML model to make a determination on whether or not the file is malicious or bengign.

## Building from Source

Binaries are provided for MacOS and Linux, Ganalyze is not tested on Windows. If building from source, ensure you have a working Go environment and run `go build` from the repository's root directory.

## Python Requirements

Ganalyze uses multiple Python libraries and has specific Python version requirements.  We recommend building a Python/Conda virtual environment to run Ganalyze. The required Python version is `3.6`. After building a virtual environment using Python 3.6, run `pip install -r requirements.txt` from the root directory of the repository to install the Python imports.

## Usage

To use the program, run `./ganalyze <directory>`, with directory being the root directory of the PE files you'd like to analyze.  Ganalyze will recursively scan through the directory and analyze each PE file. A `report` directory is produced which contains an index.html file and html files for each PE file analyzed. Index.html contains a listing of all the files analyzed and links to each of the respective html reports.  The reports directory will maintain the directory structure of the original directory passed in to Ganalyze.
