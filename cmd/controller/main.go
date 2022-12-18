package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/api/http"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/api/http/controller"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/proxyclient/queue"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/repository/mysql"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/usecase"
	"github.com/rs/zerolog/log"
)

func main() {
	port := flag.String("port", "8080", "listen port")
	dsn := flag.String("mysql", "", "mysql data source name")
	interval := flag.Int("interval", 10, "queue loop interval seconds")

	flag.Parse()

	db, err := mysql.NewDB(*dsn)
	if err != nil {
		log.Fatal().Str("dsn", *dsn).Msg("failed to connect db")
	}
	defer db.Close()
	repo := mysql.NewRepository(db)

	qu := queue.NewQueue(nil, time.Duration(*interval)*time.Second)
	go qu.Start()
	defer qu.Close()
	client := queue.NewClient(qu)

	ucProxy := usecase.NewProxy(repo, client)
	ctrlProxy := controller.NewProxy(ucProxy)

	ucRoute := usecase.NewRoute(repo, client)
	ctrlRoute := controller.NewRoute(ucRoute)

	server := http.NewServer(ctrlProxy, ctrlRoute)
	server.SetRoute()

	go server.Run(":" + *port)
	defer server.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
