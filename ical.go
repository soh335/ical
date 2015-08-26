package ical

// https://www.ietf.org/rfc/rfc2445.txt

import (
	"bufio"
	"io"
	"time"
)

const (
	stampLayout    = "20060102T150405Z"
	dateLayout     = "20060102"
	dateTimeLayout = "20060102T150405"
)

type VCalendar struct {
	PRODID   string
	VERSION  string
	CALSCALE string

	X_WR_CALNAME  string
	X_WR_TIMEZONE string
	X_WR_CALDESC  string

	VComponent []VComponent
}

func NewBasicVCalendar() *VCalendar {
	return &VCalendar{
		VERSION:  "2.0",
		CALSCALE: "GREGORIAN",
	}
}

func (c *VCalendar) Encode(w io.Writer) error {
	var b = bufio.NewWriter(w)

	if _, err := b.WriteString("BEGIN:VCALENDAR\r\n"); err != nil {
		return err
	}

	// use a slice map to preserve order during for range
	attrs := []map[string]string{
		{"VERSION:": c.VERSION},
		{"PRODID:": c.PRODID},
		{"X-WR-CALNAME:": c.X_WR_CALNAME},
		{"X-WR-CALDESC:": c.X_WR_CALDESC},
		{"X-WR-TIMEZONE:": c.X_WR_TIMEZONE},
		{"CALSCALE:": c.CALSCALE},
	}

	for _, item := range attrs {
		for k, v := range item {
			if len(v) == 0 {
				continue
			}
			if _, err := b.WriteString(k + v + "\r\n"); err != nil {
				return err
			}
		}
	}

	for _, component := range c.VComponent {
		if err := component.EncodeIcal(b); err != nil {
			return err
		}
	}

	if _, err := b.WriteString("END:VCALENDAR\r\n"); err != nil {
		return err
	}

	return b.Flush()
}

type VComponent interface {
	EncodeIcal(w io.Writer) error
}

type VEvent struct {
	UID     string
	DTSTAMP time.Time
	DTSTART time.Time
	DTEND   time.Time
	SUMMARY string
	TZID    string

	AllDay bool
}

func (e *VEvent) EncodeIcal(w io.Writer) error {

	var timeStampLayout, timeStampType string

	if e.AllDay {
		timeStampLayout = dateLayout
		timeStampType = "DATE"
	} else {
		timeStampLayout = dateTimeLayout
		timeStampType = "DATE-TIME"
	}

	b := bufio.NewWriter(w)
	if _, err := b.WriteString("BEGIN:VEVENT\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("DTSTAMP:" + e.DTSTAMP.UTC().Format(stampLayout) + "\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("UID:" + e.UID + "\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("TZID:" + e.TZID + "\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("SUMMARY:" + e.SUMMARY + "\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("DTSTART;TZID=" + e.TZID + ";VALUE=" + timeStampType + ":" + e.DTSTART.Format(timeStampLayout) + "\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("DTEND;TZID=" + e.TZID + ";VALUE=" + timeStampType + ":" + e.DTEND.Format(timeStampLayout) + "\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("END:VEVENT\r\n"); err != nil {
		return err
	}

	return b.Flush()
}
