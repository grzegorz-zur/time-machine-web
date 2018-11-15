package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type TimeMachine struct {
	PID string
}

var (
	tmDir     = regexp.MustCompile("timemachine-([0-9]+)")
	indexTmpl = template.Must(template.ParseFiles("index.html"))
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	err := http.ListenAndServe(":9876", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tms := find()
	indexTmpl.Execute(w, tms)
}

func find() (tms []TimeMachine) {
	dir := os.TempDir()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return
	}
	for _, file := range files {
		name := file.Name()
		if pids := tmDir.FindStringSubmatch(name); pids != nil {
			pid := pids[1]
			tm := TimeMachine{
				PID: pid,
			}
			tms = append(tms, tm)
		}
	}
	return tms
}
