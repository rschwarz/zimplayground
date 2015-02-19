package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	resultsDir = "results"

	modelFilename    = "model.zpl"
	solutionFilename = "scip.sol"
	outputFilename   = "output.log"

	inputTemplate  = "html/input.html"
	resultTemplate = "html/result.html"
)

func solve(dir string) (err error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return err
	}

	// TODO: actually solve
	return nil
}

func inputHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(inputTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	model := r.FormValue("model")
	fmt.Println("Model:", model)
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(model)))
	fmt.Println("Hash:", hash)

	dir := path.Join(resultsDir, hash)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filename := path.Join(dir, modelFilename)
		err = ioutil.WriteFile(filename, []byte(model), 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = solve(dir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/result/"+hash, http.StatusFound)
}

type Result struct {
	Hash     string
	Model    string
	Solution string
	Output   string
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/result/"):]

	dir := path.Join(resultsDir, hash)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &Result{Hash: hash, Model: "", Solution: "", Output: ""}

	model := path.Join(dir, modelFilename)
	if _, err := os.Stat(model); err == nil {
		content, err := ioutil.ReadFile(model)
		if err == nil {
			res.Model = string(content)
		}
	}

	sol := path.Join(dir, solutionFilename)
	if _, err := os.Stat(sol); err == nil {
		content, err := ioutil.ReadFile(sol)
		if err == nil {
			res.Solution = string(content)
		}
	}

	out := path.Join(dir, outputFilename)
	if _, err := os.Stat(out); err == nil {
		content, err := ioutil.ReadFile(out)
		if err == nil {
			res.Output = string(content)
		}
	}

	t, err := template.ParseFiles(resultTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", inputHandler)
	http.HandleFunc("/solve/", solveHandler)
	http.HandleFunc("/result/", resultHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
