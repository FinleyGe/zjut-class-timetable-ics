package generate

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Term      string `json:"term"`
	Year      string `json:"year"`
	StartDate string `json:"start_date"`
}

func GetClassTable(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	fmt.Println(req)

	classes := GetClasses(
		User{
			Username: req.Username,
			Password: req.Password,
		},
		GetClass{
			Term: req.Term,
			Year: req.Year,
		},
	)
	startTime, _ := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
	s := GenerateICS(classes, startTime)
	c.JSON(200,
		s,
	)
}
