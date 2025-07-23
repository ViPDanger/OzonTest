package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ViPDanger/OzonTest/internal/app"
)

func init() {

}

// @title           MockXMLDaily
// @version         1.0
// @description     Мок-сервис получения курсов валют с сайта cbr.ru в XML-формате.
// @host      127.0.0.1:8080
// @grpcHost  127.0.0.1:2020
// @accept    xml
// @produce   xml

var host = "127.0.0.1:8080"
var mongoURI = "mongodb://127.0.0.1:27017"
var mongodbName = "xml_daily"
var mongoUser = "admin"
var mongoPassword = "admin"
var grpcHost = "127.0.0.1:2020"

func main() {
	for i := range os.Args {
		checkArg(os.Args[i])
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGALRM)
	defer func() {
		cancel()
		time.Sleep(2 * time.Second)
		fmt.Println("Greatefull shutdown")
	}()
	_, err := app.Run(ctx, host, mongoURI, mongodbName, mongoUser, mongoPassword, grpcHost)
	if err != nil {
		fmt.Printf("main(): error to run app.Run() %v\n", err.Error())
		return
	}

	<-ctx.Done()
}

func checkArg(arg string) {
	if !strings.HasPrefix(arg, "--") {
		return
	}

	arg = strings.TrimPrefix(arg, "--")
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return
	}
	key := parts[0]
	value := parts[1]

	switch key {
	case "host":
		host = value
	case "mongoURI":
		mongoURI = value
	case "mongodbName":
		mongodbName = value
	case "mongoUser":
		mongoUser = value
	case "mongoPassword":
		mongoPassword = value
	case "grpcHost":
		grpcHost = value
	}
}
