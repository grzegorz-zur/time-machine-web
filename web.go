package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	Date    string
	Time    string
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
	then := time.Unix(time.Now().UTC().Unix()+int64(tm.Value), 0)
	vs.Date, vs.Time = then.Format("2006-01-02"), then.Format("15:04")
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
	case "datetime":
		date := r.FormValue("date")
		time := r.FormValue("time")
		value = offset(timeParse(date, time))
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

func timeParse(dt, tm string) (t time.Time) {
	text := dt + " " + tm
	layout := "2006-01-02 15:04"
	location := time.Now().Location()
	t, err := time.ParseInLocation(layout, text, location)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func offset(then time.Time) (seconds int) {
	now := time.Now()
	diff := then.Unix() - now.Unix()
	seconds = int(diff)
	return
}
