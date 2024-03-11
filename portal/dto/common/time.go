package common

import (
	"time"
)

type Time struct {
	UpTime     int  `xml:"UpTime,omitempty"`
	SummerTime int  `xml:"SummerTime,omitempty"`
	TimeZone   *int `xml:"TimeZone,omitempty"`
	Timestamp  Timestamp
}

func NewTime() *Time {
	result := &Time{}

	return result
}

func (t *Time) Now() {
	t.Timestamp = Timestamp(time.Now())
}

func (t *Time) WithoutTimeZone() {
	t.TimeZone = nil
}

func (t *Time) WithTimeZone(tz int) {
	t.TimeZone = &tz
}
