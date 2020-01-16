# Indexing Your Hugo Site with Python

This is a simple python tool for indexing a hugo site so that one can make a 
docsearch for his/her static site. 

## Dependency

One need to install the following module or make them avalibale for python 
interpreter:

1.  toml
2.  yaml
3.  json

## Usage

One should run the script `indexGen.py` in his/her hugo root directory with the 
parameter `$config.ext`, filename of the config file.

## TODO

* [x] Support for toml format config file
* [ ] Support for json/yaml format config file
* [ ] Upload the index file to Algolia automatically.

## Acknowledge

Thanks for the developer of Hugo-algolia and the owner of the script: 
[fnurl/docsearch-pageindexer.py](https://gist.github.com/fnurl/586dbdb7d313f1911580ae873d5ad213)