package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/roylic/gofolder/api"
)

func main() {

	// config
	listenAddr := flag.String("listenaddr", ":9990", "todo")
	flag.Parse()
	fmt.Println("Now start listening...")

	// handle
	http.HandleFunc("/user", api.HandleUserAPI)

	// server start
	http.ListenAndServe(*listenAddr, nil)
}
