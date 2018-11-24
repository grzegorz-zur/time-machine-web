package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	tmsTmpl = template.Must(template.ParseFiles("tms.html"))
	tmTmpl  = template.Must(template.ParseFiles("tm.html"))
)

type Values struct {
	Pid     string
	Command string
	Method  string
	Value   int
	Sign    int
	Weeks   int
	Days    int
	Hours   int
	Minutes int
	Seconds int
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
	tmsTmpl.Execute(w, list())
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	pid := strings.Trim(r.URL.Path, "/")
	tm := get(pid)
	vs := Values{
		Pid:     tm.Pid,
		Command: tm.Command,
		Value:   tm.Value,
	}
	methods := r.URL.Query()["method"]
	if len(methods) == 1 {
		vs.Method = methods[0]
	}
	vs.Sign, vs.Weeks, vs.Days, vs.Hours, vs.Minutes, vs.Seconds = split(tm.Value)
	tmTmpl.Execute(w, vs)
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	pid := strings.Trim(r.URL.Path, "/")
	method := r.FormValue("method")
	value := 0
	switch method {
	case "raw":
		value = formInt(r, "value")
	case "units":
		sign := formInt(r, "sign")
		weeks := formInt(r, "weeks")
		days := formInt(r, "days")
		hours := formInt(r, "hours")
		minutes := formInt(r, "minutes")
		seconds := formInt(r, "seconds")
		value = join(sign, weeks, days, hours, minutes, seconds)
	}
	set(pid, value)
	http.Redirect(w, r, r.URL.Path+"?method="+method, http.StatusSeeOther)
}

func formInt(r *http.Request, name string) (value int) {
	text := r.FormValue(name)
	if text == "" {
		return
	}
	v, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	value = int(v)
	return
}
