package utils

import (
	"testing"
	"time"
)

const timeFormat = "Mon Jan 2 15:04:05.000"

func testRelTime(t *testing.T, t1, t2 time.Time, expectedResult string) {
	result := RelTime(t1, t2)
	if result != expectedResult {
		t.Error("Got '" + result + "', expected '" + expectedResult + "'")
	}

	var diff time.Duration
	if t1.After(t2) {
		diff = t1.Sub(t2)
	} else {
		diff = t2.Sub(t1)
	}

	result = DurationString(diff)
	if result != expectedResult {
		t.Error("Got '" + result + "', expected '" + expectedResult + "'")
	}
}

func testCustomRelTime(t *testing.T, t1, t2 time.Time, numParts int, glue string, expectedResult string) {
	result := CustomRelTime(t1, t2, numParts, glue)
	if result != expectedResult {
		t.Error("Got '" + result + "', expected '" + expectedResult + "'")
	}

	var diff time.Duration
	if t1.After(t2) {
		diff = t1.Sub(t2)
	} else {
		diff = t2.Sub(t1)
	}
	result = CustomDurationString(diff, numParts, glue)
	if result != expectedResult {
		t.Error("Got '" + result + "', expected '" + expectedResult + "'")
	}
}

func TestRelTime(t *testing.T) {
	var t1 time.Time
	var t2 time.Time

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:00:05.000")
	testRelTime(t, t1, t2, "4 minutes")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:00:04.000")
	testRelTime(t, t1, t2, "4 minutes 1 second")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:00:03.000")
	testRelTime(t, t1, t2, "4 minutes 2 seconds")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:03.000")
	testRelTime(t, t1, t2, "2 seconds")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:04.000")
	testRelTime(t, t1, t2, "1 second")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:04.001")
	testRelTime(t, t1, t2, "now")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:04.999")
	testRelTime(t, t1, t2, "now")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:04.999")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	testRelTime(t, t1, t2, "now")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:00:03.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	testRelTime(t, t1, t2, "4 minutes 2 seconds")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 2 15:04:05.000")
	testRelTime(t, t1, t2, "1 day")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 1 15:04:05.000")
	testRelTime(t, t1, t2, "2 days")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 14:04:05.000")
	testRelTime(t, t1, t2, "1 hour")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 13:04:05.000")
	testRelTime(t, t1, t2, "2 hours")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 10 15:04:05.000")
	testRelTime(t, t1, t2, "1 week")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 17 15:04:05.000")
	testRelTime(t, t1, t2, "2 weeks")

	t1, _ = time.Parse(timeFormat, "Mon Sep 1 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Oct 1 15:04:05.000")
	testRelTime(t, t1, t2, "1 month")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Mar 3 15:04:05.000")
	testRelTime(t, t1, t2, "2 months")
}

func TestCustomRelTime(t *testing.T) {
	var t1 time.Time
	var t2 time.Time

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:00:05.000")
	testCustomRelTime(t, t1, t2, 1, " ", "4 minutes")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:00:02.000")
	testCustomRelTime(t, t1, t2, 1, " ", "4 minutes")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:00:02.000")
	testCustomRelTime(t, t1, t2, 2, " ", "4 minutes 3 seconds")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:00:02.000")
	testCustomRelTime(t, t1, t2, 3, " ", "4 minutes 3 seconds")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 3 15:00:02.001")
	testCustomRelTime(t, t1, t2, 3, " ", "4 minutes 2 seconds")

	t1, _ = time.Parse(timeFormat, "Mon Jan 3 15:04:05.000")
	t2, _ = time.Parse(timeFormat, "Mon Jan 2 14:03:04.000")
	testCustomRelTime(t, t1, t2, 1, " ", "1 day")
	testCustomRelTime(t, t1, t2, 2, " ", "1 day 1 hour")
	testCustomRelTime(t, t1, t2, 3, " ", "1 day 1 hour 1 minute")
	testCustomRelTime(t, t1, t2, 4, " ", "1 day 1 hour 1 minute 1 second")

	testCustomRelTime(t, t1, t2, 1, ", ", "1 day")
	testCustomRelTime(t, t1, t2, 2, ", ", "1 day, 1 hour")
	testCustomRelTime(t, t1, t2, 3, ", ", "1 day, 1 hour, 1 minute")
	testCustomRelTime(t, t1, t2, 4, ", ", "1 day, 1 hour, 1 minute, 1 second")
}
