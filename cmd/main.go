package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ViPDanger/OzonTest/internal/app"
)

func init() {

}

// @title           OzonTest
// @version         1.0
// @description     Мок-сервис получения курсов валют с сайта cbr.ru в XML-формате.
// @host      127.0.0.1:8080
// @accept    xml
// @produce   xml

const host = ":8080"
const mongoURI = "mongodb://127.0.0.1:27017"
const mongodbName = "xml_daily"
const mongoUser = "admin"
const mongoPassword = "admin"
const grpcAdress = ":2020"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGALRM)
	defer cancel()
	if _, err := app.Run(ctx, host, mongoURI, mongodbName, mongoUser, mongoPassword, grpcAdress); err != nil {
		cancel()
		fmt.Printf("main(): error to run app.Run() %v", err.Error())
		return
	}
	<-ctx.Done()
}
