package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/yuin/gopher-lua"
)

var jobs []*gocron.Job

// This is how the timer interface will work
// Time Parsing:
//  - space-delimited strings.
//  - Day selector:
//      - mon, monday (case-insensitive)
//      - weekday, weekend, day
//  - Time selector:
//      - 10:30, 10:30am, 10:30pm, 20:30 (HH:MM, optional seconds)

// day month year
var _timeformats = []string{
	"2/1/2006 15:04:05",
	"2/1/2006",
	"2/1",
	"2/1 15:04:05",
	"15:04:05",
	"15:04",
	"15",
	"2/1/2006 3:04:05PM",
	"2/1 3:04:05PM",
	"3:04:05PM",
	"3:04PM",
	"3PM",
}
var _weekdays = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
}
var _weekends = []time.Weekday{
	time.Saturday,
	time.Sunday,
}

// base object for storing schedulable units
// value of -1 means "every"
type schedule struct {
	years     []int
	months    []time.Month
	days      []int
	weekdays  []time.Weekday
	times     []string
	durations []*time.Duration
}

func ScheduleTask(s string, task func()) {
	// parse
	terms := strings.Split(s, " ")

	var (
		sched    = new(schedule)
		parseErr error
	)

	for _, term := range terms {
		if parseErr = parseTerm(term, sched); parseErr != nil {
			fmt.Println("PARSE ERROR", parseErr, term)
		}
	}
	sched.Schedule(task)
}

func (sched *schedule) Schedule(task func()) {
	// try durations
	if len(sched.durations) > 0 {
		for _, dur := range sched.durations {
			dur := dur
			go func() {
				fmt.Println("Scheduling Duration", dur)
				for _ = range time.Tick(*dur) {
					task()
				}
			}()
		}
		return
	}
	// check for times-only
	if len(sched.weekdays) == 0 {
		for _, tstr := range sched.times {
			fmt.Println("Scheduling Daily Timer At", tstr)
			j := gocron.Every(1).Day().At(tstr)
			j.Do(task)
			jobs = append(jobs, j)
		}
		return
	}
	// try for weekdays
	for _, weekday := range sched.weekdays {
		for _, tstr := range sched.times {
			job := gocron.Every(1)
			switch weekday {
			case time.Monday:
				job = job.Monday()
			case time.Tuesday:
				job = job.Tuesday()
			case time.Wednesday:
				job = job.Wednesday()
			case time.Thursday:
				job = job.Thursday()
			case time.Friday:
				job = job.Friday()
			case time.Saturday:
				job = job.Saturday()
			case time.Sunday:
				job = job.Sunday()
			}
			fmt.Println("Scheduling Weekly timer At", weekday, tstr)
			job = job.At(tstr)
			jobs = append(jobs, job)
			job.Do(task)
		}
	}

}

func parseTerm(term string, sched *schedule) error {
	term = strings.ToLower(term)

	// match day expressions
	var matchDayName = false
	switch term {
	case "mon", "monday":
		sched.weekdays = append(sched.weekdays, time.Monday)
		matchDayName = true
	case "tue", "tues", "tuesday":
		sched.weekdays = append(sched.weekdays, time.Tuesday)
		matchDayName = true
	case "wed", "wednesday":
		sched.weekdays = append(sched.weekdays, time.Wednesday)
		matchDayName = true
	case "thu", "thur", "thurs", "thursday":
		sched.weekdays = append(sched.weekdays, time.Thursday)
		matchDayName = true
	case "fri", "friday":
		sched.weekdays = append(sched.weekdays, time.Friday)
		matchDayName = true
	case "sat", "saturday":
		sched.weekdays = append(sched.weekdays, time.Saturday)
		matchDayName = true
	case "sun", "sunday":
		sched.weekdays = append(sched.weekdays, time.Sunday)
		matchDayName = true
	case "weekday":
		sched.weekdays = append(sched.weekdays, _weekdays...)
		matchDayName = true
	case "weekend":
		sched.weekdays = append(sched.weekdays, _weekends...)
		matchDayName = true
	case "day":
		sched.weekdays = append(sched.weekdays, _weekdays...)
		sched.weekdays = append(sched.weekdays, _weekends...)
		matchDayName = true
	}

	if matchDayName {
		return nil
	}

	// try parse duration
	dur, err := ParseDuration(term)
	if err == nil {
		sched.durations = append(sched.durations, dur)
		return nil
	}

	sched.times = append(sched.times, term)

	// match explicit time formats
	var matchedTime time.Time
	for _, layout := range _timeformats {
		t, err := time.Parse(layout, term)
		if err != nil {
			if _, ok := err.(*time.ParseError); ok {
			} else {
				fmt.Println("parse error unk", err)
				return err
			}
		} else {
			matchedTime = t
			break
		}
	}
	year, month, day := matchedTime.Date()
	if year != 0 {
		sched.years = append(sched.years, year)
	}
	if month != 0 {
		sched.months = append(sched.months, month)
	}
	if day != 0 {
		sched.days = append(sched.days, day)
	}

	return nil
}

func startCronScheduler(L *lua.LState) {
	go func() {
		for _ = range time.Tick(5 * time.Second) {
			_, t := gocron.NextRun()
			fmt.Println(t, time.Now())
		}
	}()
	go func() {
		<-gocron.Start()
	}()
}
