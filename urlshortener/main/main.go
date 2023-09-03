package main

import (
	"fmt"
	"net/http"
	"flag"
	"os"

	"github.com/hamzanaciri99/golang-ex/urlshortener"
	"github.com/hamzanaciri99/golang-ex/util"
)

//Define and initialize flags
var (
	file string
)
func flagsInit() {
	flag.StringVar(&file, "file", "", "yaml file containing path/url mapping")
}

func init() {
	flagsInit()
	flag.Parse()
}	

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := os.ReadFile("urlshortener/main/routes.yaml")
	util.CheckError(err)

	yamlHandler, err := urlshortener.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}