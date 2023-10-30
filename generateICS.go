package main

import (
	// "fmt"

	"fmt"

	ics "github.com/arran4/golang-ical"
)

func generateICS(classes []Class) string {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetVersion("2.0")
	cal.SetXWRCalName("浙江工业大学课程表")
	// cal.SetTimezoneId("Asia/Shanghai")
	for _, class := range classes {
		start, end, repeat := class.getTime()
		// for r := 0; r < repeat; r++ {
		lesson := cal.AddEvent(class.getUID())
		lesson.SetSummary(class.LessonName)
		lesson.SetDescription(class.TeacherName)
		lesson.SetLocation(class.LessonPlace)
		lesson.SetStartAt(start)
		lesson.SetEndAt(end)
		// fmt.Println(lesson.Serialize())
		// }
		if repeat > 1 {
			lesson.AddRrule(fmt.Sprintf("FREQ=WEEKLY;COUNT=%d", repeat))
		}
	}

	return cal.Serialize()
}
