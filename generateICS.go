package main

import (
	"log"
	"os"
)

type ICS struct {
	ProdID   string
	Version  string
	CalScale string
	Method   string
	Events   []Event
}

type Event struct {
	UID         string
	DTSTAMP     string
	DTSTART     string
	DTEND       string
	SUMMARY     string
	DESCRIPTION string
	LOCATION    string
}

func GenerateICSFile(ics ICS, filename string) {
	// generate ics file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	file.WriteString("BEGIN:VCALENDAR\n")
	file.WriteString("PRODID:" + ics.ProdID + "\n")
	file.WriteString("VERSION:" + ics.Version + "\n")
	file.WriteString("CALSCALE:" + ics.CalScale + "\n")
	file.WriteString("METHOD:" + ics.Method + "\n")
	for _, event := range ics.Events {
		file.WriteString("BEGIN:VEVENT\n")
		file.WriteString("UID:" + event.UID + "\n")
		file.WriteString("DTSTAMP:" + event.DTSTAMP + "\n")
		file.WriteString("DTSTART:" + event.DTSTART + "\n")
		file.WriteString("DTEND:" + event.DTEND + "\n")
		file.WriteString("SUMMARY:" + event.SUMMARY + "\n")
		file.WriteString("DESCRIPTION:" + event.DESCRIPTION + "\n")
		file.WriteString("LOCATION:" + event.LOCATION + "\n")
		file.WriteString("END:VEVENT\n")
	}
	file.WriteString("END:VCALENDAR\n")
}
