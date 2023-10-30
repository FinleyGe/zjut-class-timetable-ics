package generate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	ics "github.com/arran4/golang-ical"
)

var classTime []string = []string{
	"08:00:00",
	"08:55:00",
	"09:55:00",
	"10:50:00",
	"11:50:00",

	"13:30:00",
	"14:25:00",
	"15:25:00",
	"16:20:00",

	"18:30:00",
	"19:25:00",
	"20:10:00",
	"21:05:00",
}

const classStep string = "00:45:00"

type Class struct {
	LessonName  string `json:"className"`
	TeacherName string `json:"teacherName"`
	LessonPlace string `json:"lessonPlace"`
	Weekday     string `json:"weekday"`
	Week        string `json:"week"`
	Type        string `json:"type"`
	Sections    string `json:"sections"`
}

func (c Class) getUID() string {
	return c.LessonName + c.TeacherName + c.LessonPlace + c.Week + c.Weekday
}

func (c Class) getTime(TermStartDate time.Time) (time.Time, time.Time, int) {
	fmt.Println(c)
	var weekday int = int(c.Weekday[0] - '0')
	var weekstart int
	var weekend int
	c.Week = c.Week[:len(c.Week)-3]

	if utf8.RuneCountInString(c.Week) > 2 {
		s := strings.Split(c.Week, "-")
		weekstart, _ = strconv.Atoi(s[0])
		weekend, _ = strconv.Atoi(s[1])
	} else {
		weekstart, _ = strconv.Atoi(c.Week)
		weekend = weekstart
	}

	var classDate time.Time = TermStartDate.AddDate(0, 0, (weekstart-1)*7+weekday-1)
	sections := strings.Split(c.Sections, "-")
	startSection, _ := strconv.Atoi(sections[0]) // the first section
	endSection, _ := strconv.Atoi(sections[1])   // the last section
	classBegin, _ := time.ParseInLocation(time.TimeOnly, classTime[startSection-1], time.Local)
	classEnd, _ := time.ParseInLocation(time.TimeOnly, classTime[endSection-1], time.Local)
	classEnd = classEnd.Add(time.Duration(time.Minute * 45))
	classBegin = classBegin.AddDate(classDate.Year(), int(classDate.Month())-1, classDate.Day()-1)
	classEnd = classEnd.AddDate(classDate.Year(), int(classDate.Month())-1, classDate.Day()-1)
	fmt.Println(classBegin, classEnd, weekend-weekstart+1)
	return classBegin, classEnd, weekend - weekstart + 1
}

func GenerateICS(classes []Class, TermStartDate time.Time) string {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetVersion("2.0")
	cal.SetXWRCalName("浙江工业大学课程表")
	// cal.SetTimezoneId("Asia/Shanghai")
	for _, class := range classes {
		start, end, repeat := class.getTime(TermStartDate)
		lesson := cal.AddEvent(class.getUID())
		lesson.SetSummary(class.LessonName)
		lesson.SetDescription(class.TeacherName)
		lesson.SetLocation(class.LessonPlace)
		lesson.SetStartAt(start)
		lesson.SetEndAt(end)
		if repeat > 1 {
			lesson.AddRrule(fmt.Sprintf("FREQ=WEEKLY;COUNT=%d", repeat))
		}
	}
	return cal.Serialize()
}
