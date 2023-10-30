package main

import (
	"zjut_class_timetable_ics/config"
	"zjut_class_timetable_ics/server"
)

func main() {
	server.Run(config.Port)
	// Login()
	// GetClasses()
	// s := generateICS(classes)
	// fmt.Println(classes)
	// f, err := os.Create("class.ics")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer f.Close()
	// f.WriteString(s)
}
