package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type TimeMachine struct {
	Pid     string
	Command string
	Value   int
}

const (
	second = 1
	minute = 60 * second
	hour   = 60 * minute
	day    = 24 * hour
	week   = 7 * day
)

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
			if tm.Command != "" {
				tms = append(tms, tm)
			}
		}
	}
	return
}

func get(pid string) (tm TimeMachine) {
	tm.Pid = pid
	tm.Command = command(tm.Pid)
	file := path.Join(os.TempDir(), "timemachine-"+tm.Pid, "get")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return
	}
	text := string(data)
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	tm.Value = int(value)
	return
}

func set(pid string, value int) {
	file := path.Join(os.TempDir(), "timemachine-"+pid, "set")
	text := strconv.Itoa(value)
	data := []byte(text)
	err := ioutil.WriteFile(file, data, os.ModeNamedPipe|0775)
	if err != nil {
		log.Println(err)
		return
	}
}

func command(pid string) (cmd string) {
	file := path.Join("/proc", pid, "cmdline")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return
	}
	cmd = strings.Replace(string(data), "\x00", " ", -1)
	return
}

func split(value int) (sign, weeks, days, hours, minutes, seconds int) {
	if value < 0 {
		sign = -1
		value = -value
	} else {
		sign = 1
	}
	weeks, value = value/week, value%week
	days, value = value/day, value%day
	hours, value = value/hour, value%hour
	minutes, value = value/minute, value%minute
	seconds = value
	return
}

func join(sign, weeks, days, hours, minutes, seconds int) (value int) {
	value += weeks * week
	value += days * day
	value += hours * hour
	value += minutes * minute
	value += seconds
	value *= sign
	return
}
