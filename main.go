package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type TimeMachine struct {
	Pid   string
	Value int64
}

var (
	pattern = regexp.MustCompile("timemachine-([0-9]+)")
	tms     = template.Must(template.ParseFiles("tms.html"))
	tm      = template.Must(template.ParseFiles("tm.html"))
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	err := http.ListenAndServe(":9876", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		listHandler(w, r)
	case r.Method == "GET":
		getHandler(w, r)
	case r.Method == "POST":
		setHandler(w, r)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	tms.Execute(w, list())
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	pid := strings.Trim(r.URL.Path, "/")
	tm.Execute(w, get(pid))
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	pid := strings.Trim(r.URL.Path, "/")
	value := r.FormValue("value")
	set(pid, value)
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func list() (tms []TimeMachine) {
	dir := os.TempDir()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return
	}
	for _, file := range files {
		name := file.Name()
		if pids := pattern.FindStringSubmatch(name); pids != nil {
			pid := pids[1]
			tm := get(pid)
			tms = append(tms, tm)
		}
	}
	return
}

func get(pid string) (tm TimeMachine) {
	tm.Pid = pid
	file := path.Join(os.TempDir(), "timemachine-"+tm.Pid, "get")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return
	}
	text := strings.TrimSpace(string(data))
	tm.Value, err = strconv.ParseInt(text, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func set(pid, value string) {
	file := path.Join(os.TempDir(), "timemachine-"+pid, "set")
	data := []byte(value)
	err := ioutil.WriteFile(file, data, os.ModeNamedPipe|0775)
	if err != nil {
		log.Println(err)
		return
	}
}
