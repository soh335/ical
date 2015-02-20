package ical

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func Test_Encode(t *testing.T) {
	c := NewBasicVCalendar()
	c.PRODID = "proid"
	c.X_WR_TIMEZONE = "Asia/Tokyo"
	c.X_WR_CALNAME = "name"
	c.X_WR_CALDESC = "desc"

	zone := time.FixedZone("Asia/Tokyo", 60*60*9)
	d := time.Date(2014, time.Month(1), 1, 0, 0, 0, 0, zone)

	c.VComponent = []VComponent{
		&VEvent{
			UID:     "123",
			DTSTAMP: d,
			DTSTART: d,
			DTEND:   d,
			SUMMARY: "summary",
			TZID:    "Asia/Tokyo",
		},
	}

	var b bytes.Buffer
	if err := c.Encode(&b); err != nil {
		t.Error("got err:", err)
	}

	expect := `BEGIN:VCALENDAR
PRODID:proid
CALSCALE:GREGORIAN
VERSION:2.0
X-WR-CALNAME:name
X-WR-CALDESC:desc
X-WR-TIMEZONE:Asia/Tokyo
BEGIN:VEVENT
DTSTAMP:20131231T150000Z
UID:123
TZID:Asia/Tokyo
SUMMARY:summary
DTSTART;TZID=Asia/Tokyo:20140101T000000
DTEND;TZID=Asia/Tokyo:20140101T000000
END:VEVENT
END:VCALENDAR
`
	expect = strings.Replace(expect, "\n", "\r\n", -1)

	if s := b.String(); s != expect {
		t.Errorf("should %v. but got %v", expect, s)
	}
}
