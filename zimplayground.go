package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"
)

const (
	resultsDir = "results"

	modelFilename    = "model.zpl"
	solutionFilename = "scip.sol"
	outputFilename   = "output.log"

	queueLength = 1024
)

var (
	timeLimitSec  = flag.Int("time", 3*60, "SCIP time limit (s)")
	memoryLimitMB = flag.Int("mem", 100, "SCIP memory limit (MB)")
	sleepTime     = flag.Int("sleep", 100, "sleep before redirect to results (ms)")
	address       = flag.String("address", ":8080", "hostname:port of server")
	processLimit  = flag.Int("processes", 4, "limit on number of SCIP processes")
	scipExec      = flag.String("scipExec", "scip", "(path to) scip executable")
)

// a job is identified by the path to run in
type Job struct {
	dir string
}

// channel-based semaphore
type Sem chan int

// the solveHandler submits jobs here, the workers get them
var queue = make(chan Job, queueLength)

// take jobs from queue and start processes
// cf. github.com/golang/go/wiki/BoundingResourceUse
func processQueue(sem Sem) {
	for {
		sem <- 1 // block until there's capacity to start a new process
		job := <-queue
		go runSolver(job, sem) // don't wait for solver to finish.
	}
}

// run subprocess and wait to finish
func runSolver(job Job, sem Sem) {
	if _, err := os.Stat(job.dir); os.IsNotExist(err) {
		return
	}

	outputFile, err := os.Create(path.Join(job.dir, outputFilename))
	if err != nil {
		return
	}
	defer outputFile.Close()

	commands := fmt.Sprintf("set limits time %d "+
		"set limits memory %d "+
		"read %s "+
		"optimize "+
		"display statistics "+
		"write solution %s "+
		"quit",
		*timeLimitSec, *memoryLimitMB,
		modelFilename, solutionFilename)
	cmd := exec.Command(*scipExec, "-c", commands)
	cmd.Dir = job.dir
	cmd.Stdout = outputFile
	cmd.Stderr = outputFile

	if err := cmd.Start(); err != nil {
		log.Printf("Solver in %s failed to start.", job.dir)
	}
	log.Printf("Solver in %s started with PID %d", job.dir, cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Solver in %s failed with code %d.",
					job.dir, status.ExitStatus())
				fmt.Fprintf(outputFile, "Terminated with return code %d!",
					status.ExitStatus())
			}
		} else {
			log.Print("Solver in %s failed: %v", job.dir, err)
		}
	}

	log.Printf("Solver finished in %s", job.dir)

	<-sem // done, resource freed
}

// submit job to queue
func submit(job Job) (err error) {
	if _, err := os.Stat(job.dir); os.IsNotExist(err) {
		return err
	}

	log.Printf("Submitting job for %s", job.dir)

	// submit job to queue, never block
	go func() {
		queue <- job
	}()

	return nil
}

type Input struct {
	Model string
}

func inputHandler(w http.ResponseWriter, r *http.Request) {
	// optional prefilled model text
	prefilled := r.FormValue("prefilled")
	input := &Result{Model: prefilled}

	err := inputTemplate.Execute(w, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func solveHandler(w http.ResponseWriter, r *http.Request) {
	model := r.FormValue("model")
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(model)))

	log.Printf("Request for model %s", hash)

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

		err = submit(Job{dir})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// add short sleep so that result might already exist when we
	// finally redirect.
	time.Sleep(time.Duration(*sleepTime) * time.Millisecond)

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
	flag.Parse()

	// check that solver exists
	scipPath, err := exec.LookPath(*scipExec)
	if err != nil {
		log.Fatalf("No solver executable found at %s!", *scipExec)
	} else {
		log.Printf("Using solver executable at %s.", scipPath)
	}

	// semaphore for number of go routines starting processes
	var sem = make(Sem, *processLimit)
	go processQueue(sem)

	http.HandleFunc("/", inputHandler)
	http.HandleFunc("/input/", inputHandler)
	http.HandleFunc("/solve/", solveHandler)
	http.HandleFunc("/result/", resultHandler)

	log.Printf("listening on %s", *address)
	log.Fatal(http.ListenAndServe(*address, nil))
}
