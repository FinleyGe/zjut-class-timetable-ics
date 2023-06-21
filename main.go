package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/viper"
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
	"16:10:00",
	"17:05:00",

	"18:30:00",
	"19:25:00",
	"20:10:00",
	"21:05:00",
}

const classStep string = "00:45:00"

var LoginUrl string
var ClassUrl string
var Cookie string
var user User
var getClass GetClass
var classes []Class
var TermStartDate time.Time

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetClass struct {
	Term string `json:"term"`
	Year string `json:"year"`
}

type classResp struct {
	Code int `json:"code"`
	Data struct {
		Info struct {
			Classname string `json:"Classname"`
			Name      string `json:"Name"`
		} `json:"info"`
		Lessons []Class `json:"lessonsTable"`
	} `json:"data"`
	Msg string `json:"msg"`
}

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

func (c Class) getTime() (time.Time, time.Time, int) {
	fmt.Println(c)
	var weekday int = int(c.Weekday[0] - '0')
	var weekstart int
	var weekend int

	// Remove the last character '周'

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

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	LoginUrl = viper.GetString("login_url")
	ClassUrl = viper.GetString("class_url")
	user.Username = viper.GetString("username")
	user.Password = viper.GetString("password")
	getClass.Term = viper.GetString("term")
	getClass.Year = viper.GetString("year")
	// set time zone to UTC+8

	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc

	TermStartDate, _ = time.ParseInLocation("2006-01-02", viper.GetString("start"), loc)
	if err != nil {
		log.Fatalln(err)
	}

	println(TermStartDate.String())
}

func Login() {
	data, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", LoginUrl, bytes.NewBuffer(data))
	resp, err := http.DefaultClient.Do(req)
	Cookie = resp.Header.Get("Set-Cookie")
	if err != nil {
		log.Fatalln(err)
	}
}

func GetClasses() {
	data, _ := json.Marshal(getClass)
	req, _ := http.NewRequest("POST", ClassUrl, bytes.NewBuffer(data))
	req.Header.Set("Cookie", Cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	var classResp classResp
	json.NewDecoder(resp.Body).Decode(&classResp)
	classes = classResp.Data.Lessons
	for i := 0; i < len(classes); i++ {
		if strings.Contains(classes[i].Week, ",") {
			weeks := strings.Split(classes[i].Week, ",")
			classes[i].Week = weeks[0]
			for j := 1; j < len(weeks); j++ {
				classes[i].Week = weeks[j]
				classes = append(classes, classes[i])
			}
		}
	}
}

func main() {
	initConfig()
	Login()
	GetClasses()
	s := generateICS(classes)
	// println(s)
	// save to file
	f, err := os.Create("class.ics")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	f.WriteString(s)
}
