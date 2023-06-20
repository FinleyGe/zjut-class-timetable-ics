package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

var classTime []string = []string{
	"08:00",
	"08:55",
	"09:55",
	"10:50",
	"11:50",

	"13:30",
	"14:25",
	"15:25",
	"16:10",
	"17:05",

	"18:30",
	"19:25",
	"20:10",
	"21:05",
}

const classStep string = "00:45"

var LoginUrl string
var ClassUrl string
var Cookie string
var user User
var getClass GetClass
var classes []class

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
		Lessons []class `json:"lessonsTable"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type class struct {
	LessonName  string `json:"className"`
	TeacherName string `json:"teacherName"`
	LessonPlace string `json:"lessonPlace"`
	Weekday     string `json:"weekday"`
	Week        string `json:"week"`
	Type        string `json:"type"`
	Sections    string `json:"sections"`
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
}

func main() {
	initConfig()
	Login()
	GetClasses()
	ics := ICS{
		ProdID:  "zjut-classtable",
		Version: "2.0",
		Events:  []Event{
			// {
			// 	UID:         "1",
			// 	DTSTAMP:     "20200101T000000Z",
			// 	DTSTART:     "20200101T000000Z",
			// 	DTEND:       "20200101T000000Z",
			// 	SUMMARY:     "test",
			// 	DESCRIPTION: "test",
			// 	LOCATION:    "test",
			// },
		},
	}
	for _, class := range classes {
		ics.Events = append(ics.Events, Event{
			UID:         class.LessonName,
			DTSTAMP:     "20200101T000000Z",
			DTSTART:     "20200101T000000Z",
			DTEND:       "20200101T000000Z",
			SUMMARY:     class.LessonName,
			DESCRIPTION: class.TeacherName,
			LOCATION:    class.LessonPlace,
		})
	}
	GenerateICSFile(ics, "test.ics")
}
