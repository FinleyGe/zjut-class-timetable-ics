package server

import (
	"log"
	"os"
	"os/signal"
	"zjut_class_timetable_ics/generate"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Term      string `json:"term"`
	Year      string `json:"year"`
	StartDate string `json:"start_date"`
}

func Run(port string) {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = append(
		corsConfig.AllowOrigins,
		"http://localhost:5173",
	)
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))
	router.POST("/", generate.GetClassTable)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		log.Println("Server is running on port " + port)
		router.Run(":" + port)
	}()

	<-signals
	log.Println("Server is shutting down...")
}
