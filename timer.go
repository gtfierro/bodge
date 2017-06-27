package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/yuin/gopher-lua"
)

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
	"2/1/2006 3:04:05 PM",
	"2/1 3:04:05 PM",
	"3:04:05 PM",
	"3:04 PM",
	"3 PM",
}

var cron *cronScheduler

const _MAXTASKS = 16384

type task struct {
}

// returns the next time this should be run
func (t *task) next() time.Time {
	return time.Now()
}

// returns true if the task is overdue to run
func (t *task) shouldRun() bool {
	return true
}

func (t *task) run() {
	fmt.Println("running")
}

func parseScheduleString(s string) {
	terms := strings.Split(s, " ")

	for _, term := range terms {
		parseTerm(term)
	}

}

func parseTerm(term string) {
	// lowercase
	term = strings.ToLower(term)

	now := time.Now()

	// match day
	var ahead = uint(0) //
	var matchDayName = false
	switch term {
	case "mon", "monday":
		ahead = uint((now.Weekday() - time.Monday) % 7)
		matchDayName = true
	case "tue", "tues", "tuesday":
		ahead = uint((now.Weekday() - time.Tuesday) % 7)
		matchDayName = true
	case "wed", "wednesday":
		ahead = uint((now.Weekday() - time.Wednesday) % 7)
		matchDayName = true
	case "thu", "thur", "thurs", "thursday":
		ahead = uint((now.Weekday() - time.Thursday) % 7)
		matchDayName = true
	case "fri", "friday":
		ahead = uint((now.Weekday() - time.Friday) % 7)
		matchDayName = true
	case "sat", "saturday":
		ahead = uint((now.Weekday() - time.Saturday) % 7)
		matchDayName = true
	case "sun", "sunday":
		ahead = uint((now.Weekday() - time.Sunday) % 7)
		matchDayName = true
	case "weekday":
		matchDayName = true
	case "weekend":
		matchDayName = true
	case "day":
		matchDayName = true
	}
	t := now.AddDate(0, 0, int(ahead))
	fmt.Println(t.Sub(now))

	if matchDayName {

		return
	}

	// match time
	var matchTime time.Time
	for _, layout := range _timeformats {
		t, err := time.Parse(layout, term)
		if err != nil {
			if pe, ok := err.(*time.ParseError); ok {
				fmt.Println("parse error", pe)
			} else {
				fmt.Println("parse error unk", err)
			}
		} else {
			matchTime = t
			fmt.Println(t)
			break
		}
	}

	y, m, d := matchTime.Date()
	fmt.Println("y,m,d", y, m, d)
}

type cronScheduler struct {
	tasks chan *task
}

func (cs *cronScheduler) addTask(t *task) {
	cs.tasks <- t
}

func (cs *cronScheduler) run() {
	for {
		time.Sleep(1)
		for t := range cs.tasks {
			if t.shouldRun() {
				t.run()
			}
		}
	}
}

func startCronScheduler(L *lua.LState) {
	cron = &cronScheduler{
		tasks: make(chan *task, _MAXTASKS),
	}
	go cron.run()

	parseScheduleString("monday 10:30")

}
