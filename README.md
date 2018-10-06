# UniDoc - Pdf optimization WebAssembly example

[UniDoc](http://unidoc.io) is a powerful PDF library for Go (golang). The library is written and supported by the owners of the [FoxyUtils.com](https://foxyutils.com) website, where the library is used to power many of the PDF services offered. 

## How to run
~~~
go run server.go
~~~ 

## Build from the source

### Install UniDoc
~~~
go get -u github.com/unidoc/unidoc/...
~~~

### Build the WebAssembly library
~~~
GOARCH=wasm GOOS=js go build -o lib.wasm wasm/main.go
~~~

