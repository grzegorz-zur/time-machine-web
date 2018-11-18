package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
)

type TimeMachine struct {
	Pid   string
	Value int64
}

var (
	pattern = regexp.MustCompile("timemachine-([0-9]+)")
)

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
	text := string(data)
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