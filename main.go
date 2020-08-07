package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

func main() {
	year, _ := strconv.Atoi(time.Now().Format("2006"))
	month, _ := strconv.Atoi(time.Now().Format("01"))

	var retval map[string]time.Duration

	// previous month
	if month == 1 {
		retval = getMonthStats(year - 1, 12)
	} else {
		retval = getMonthStats(year, month-1)
	}
	fmt.Println("elozo honap")
	for user, dur := range retval {
		fmt.Printf("%s: %s\n", user, dur.String())
	}
	fmt.Println()

	//current month
	retval = getMonthStats(year, month)
	fmt.Println("aktualis honap")
	for user, dur := range retval {
		fmt.Printf("%s: %s\n", user, dur.String())
	}
	fmt.Println()

	//next month
	if month == 12 {
		retval = getMonthStats(year+1, 1)
	} else {
		retval = getMonthStats(year, month+1)
	}
	fmt.Println("kovetkezo honap")
	for user, dur := range retval {
		fmt.Printf("%s: %s\n", user, dur.String())
	}

}

func getMonthStats(year, month int) map[string]time.Duration {
	pd := pagerduty.NewClient(pagerDutyToken)

	schedules, err := pd.ListSchedules(pagerduty.ListSchedulesOptions{})
	if err != nil {
		panic(err)
	}

	since := fmt.Sprintf("%d%.2d01-000000", year, month)
	until := fmt.Sprintf("%d%.2d01-000000", year, month+1)
	if month == 12 {
		until = fmt.Sprintf("%d%.2d01-000000", year+1, 1)
	}

	o := pagerduty.GetScheduleOptions{
		Since:    since,
		Until:    until,
		TimeZone: "Europe/Budapest",
	}
	schedule, err := pd.GetSchedule(schedules.Schedules[0].ID, o)
	if err != nil {
		panic(err)
	}

	onCall := make(map[string]time.Duration)

	for _, entry := range schedule.FinalSchedule.RenderedScheduleEntries {
		startTime, _ := time.Parse(time.RFC3339, entry.Start)
		endTime, _ := time.Parse(time.RFC3339, entry.End)

		dur := endTime.Sub(startTime)
		onCall[entry.User.Summary] += dur
	}

	return onCall
}
