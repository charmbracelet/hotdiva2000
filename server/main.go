package main

import (
	"io"
	"net/http"

	"github.com/charmbracelet/hotdiva2000"
	"github.com/charmbracelet/log"
)

func generateName(w http.ResponseWriter, r *http.Request) {
	name := hotdiva2000.Generate()
	log.Infof("Generated name: %s", name)
	_, err := io.WriteString(w, name)
	if err != nil {
		log.Error(err)
	}
}

func main() {
	http.HandleFunc("/", generateName)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}
