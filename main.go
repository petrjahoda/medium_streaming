package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"net/http"
	"os"
)

const serviceName = "Web application login"
const serviceDescription = "Web application login"
const connection = "user=postgres password=password dbname=medium host=localhost port=5433 sslmode=disable"
const databaseConnection = "user=postgres password=password dbname=postgres host=localhost port=5433 sslmode=disable"

type program struct{}

func main() {
	fmt.Println(serviceName + " starting...")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println("Cannot start: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		fmt.Println("Cannot start: " + err.Error())
	}
}

func (p *program) Start(service.Service) error {
	fmt.Println(serviceName + " started")
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	fmt.Println(serviceName + " stopped")
	return nil
}

func (p *program) run() {
	go checkDatabase()
	router := httprouter.New()
	sseStreamer := sse.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.GET("/", homepage)
	router.Handler("GET", "/data", sseStreamer)
	go insertRandomDataToDatabase()
	go streamDataToWebPage(sseStreamer)
	err := http.ListenAndServe(":80", router)
	if err != nil {
		fmt.Println("Problem starting service: " + err.Error())
		os.Exit(-1)
	}
	fmt.Println(serviceName + " running")
}
