package main

import (
	"strings"
	"time"
)

type timeRange []time.Time

func (tr timeRange) Len() int {
	return len(tr)
}

func (tr timeRange) Swap(i, j int) {
	tr[i], tr[j] = tr[j], tr[i]
}

func (tr timeRange) Less(i, j int) bool {
	return tr[i].Before(tr[j])
}

func formatTimeSpan(timeSpan string) (timeRange, error) {
	//Mon Jan 2 15:04:05 -0700 MST 2006 ... "02/01/2006 15:04:05.000 -0700"
	var tr timeRange

	ts := strings.Split(timeSpan, "-")

	for _, v := range ts {
		t, err := time.Parse("02/01/2006 15:04:05.000 -0700", v)
		if err != nil {
			return nil, err
		}

		tr = append(tr, t)
	}
	return tr, nil
}

func (tr timeRange) timeRangeToString() []string {
	var s []string
	for _, v := range tr {
		s = append(s, timeToString(v))
	}
	return s
}

func timeToString(t time.Time) string {
	return t.Format("02/01/2006 15:04:00")
}

func stringToTime(s string) (time.Time, error) {
	t, err := time.Parse("02/01/2006 15:04:05.000 -0700", s)
	if err != nil {
		return t, err
	}
	return t, nil
}
