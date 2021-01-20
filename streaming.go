package main

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/sse"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

func streamDataToWebPage(sseStreamer *sse.Streamer) {
	_, err := sql.Open("postgres", databaseConnection)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	listener := pq.NewListener(databaseConnection, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		waitForNotification(listener, sseStreamer)
	}
}

func waitForNotification(listener *pq.Listener, sseStreamer *sse.Streamer) {
	for {
		select {
		case n := <-listener.Notify:
			sseStreamer.SendString("data", "data", n.Extra)
			return
		}
	}
}

func insertRandomDataToDatabase() {
	for {
		productionDatabase, err := gorm.Open(postgres.Open(databaseConnection), &gorm.Config{})
		productionDB, _ := productionDatabase.DB()
		if err != nil {
			fmt.Println("Problem opening database, looks like it does not exist")
			_ = productionDB.Close()
			time.Sleep(2 * time.Second)
			continue
		}
		randomNumber := strconv.Itoa(rand.Intn(100-0) + 0)
		fmt.Println("Inserting random data to database: " + randomNumber)
		data := Data{Data: randomNumber}
		productionDatabase.Save(&data)
		_ = productionDB.Close()
		time.Sleep(2 * time.Second)
	}
}
