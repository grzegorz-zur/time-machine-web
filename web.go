package main

import (
	"html/template"
	"net/http"
	"strings"
)

var (
	tms = template.Must(template.ParseFiles("tms.html"))
	tm  = template.Must(template.ParseFiles("tm.html"))
)

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
