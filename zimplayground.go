package main

import (
	"html/template"
	"net/http"
)

var inputTemplate = "html/input.html"

func inputHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(inputTemplate)
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", inputHandler)
	http.ListenAndServe(":8080", nil)
}
