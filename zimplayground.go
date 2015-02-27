package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

const (
	resultsDir = "results"

	modelFilename    = "model.zpl"
	solutionFilename = "scip.sol"
	outputFilename   = "output.log"

	timeLimitSec  = 3 * 60
	memoryLimitMB = 100
)

func runSolver(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return
	}

	// TODO: add limit on number of parallel solver runs

	commands := fmt.Sprintf("set limits time %d "+
		"set limits memory %d "+
		"read %s  opt  write solution %s  quit",
		timeLimitSec, memoryLimitMB,
		modelFilename, solutionFilename)
	cmd := exec.Command("scip", "-c", commands, "-l", outputFilename)
	cmd.Dir = dir
	_ = cmd.Run()
}

func solve(dir string) (err error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return err
	}

	// for now, just start new process for every call
	go runSolver(dir)

	return nil
}

func inputHandler(w http.ResponseWriter, r *http.Request) {
	err := inputTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	model := r.FormValue("model")
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(model)))

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

	// add short sleep so that result might already exist when we
	// finally redirect.
	time.Sleep(100 * time.Millisecond)

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

	if err := resultTemplate.Execute(w, res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", inputHandler)
	http.HandleFunc("/solve/", solveHandler)
	http.HandleFunc("/result/", resultHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
