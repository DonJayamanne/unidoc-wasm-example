package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	listen := "127.0.0.1:8080"
	fmt.Printf("Started at http://%s\n", listen)
	log.Fatalln(http.ListenAndServe(listen, http.FileServer(http.Dir("."))))
}
