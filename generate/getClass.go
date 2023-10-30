package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"zjut_class_timetable_ics/config"
)

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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetClass struct {
	Term string `json:"term"`
	Year string `json:"year"`
}

func login(user User) string {
	data, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", config.LoginUrl, bytes.NewBuffer(data))
	resp, err := http.DefaultClient.Do(req)
	Cookie := resp.Header.Get("Set-Cookie")
	if err != nil {
		log.Fatalln(err)
	}
	return Cookie
}

func GetClasses(user User, getClass GetClass) []Class {
	Cookie := login(user)
	data, _ := json.Marshal(getClass)
	req, _ := http.NewRequest("POST", config.ClassUrl, bytes.NewBuffer(data))
	req.Header.Set("Cookie", Cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	var classResp classResp
	json.NewDecoder(resp.Body).Decode(&classResp)

	classes := classResp.Data.Lessons
	for i := 0; i < len(classes); i++ {
		if strings.Contains(classes[i].Week, ",") {
			weeks := strings.Split(classes[i].Week, ",")
			classes[i].Week = weeks[0]
			for j := 0; j < len(weeks); j++ {
				classes[i].Week = weeks[j]
				// classes[i].LessonName = classes[i].LessonName + "-" + strconv.Itoa(j+1)
				classes = append(classes, classes[i])
			}
		}
	}
	fmt.Println(classes)
	return classes
}
