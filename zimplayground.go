package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var inputTemplate = "html/input.html"

func inputHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(inputTemplate)
	t.Execute(w, nil)
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	model := r.FormValue("model")
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(model)))
	// TODO: solve problem

	http.Redirect(w, r, "/result/"+hash, http.StatusFound)
}

func main() {
	http.HandleFunc("/", inputHandler)
	http.HandleFunc("/solve", solveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
