package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// получаем переменные окружения
	socket := os.Getenv("SERVER_LISTEN_SOCKET")
	if socket == "" {
		log.Fatal("environment variable SERVER_LISTEN_SOCKET must be set")
	}
	connStringPostgres := os.Getenv("POSTGRES_CONN_STRING")
	if connStringPostgres == "" {
		log.Fatal("environment variable POSTGRES_CONN_STRING must be set")
	}

	// создаем образ БД
	var bd storage.Model
	bd, err := postgres.New(connStringPostgres)
	if err != nil {
		log.Fatalf("error connecting to database [%v]\n", err)
	}
	defer bd.Close()

	// создаем API сервера
	l := log.New(os.Stderr, "[GoNews server]\t->\t", log.LstdFlags|log.Lmsgprefix)
	api := api.New(bd, l)

	// конфигурируем сервер
	srv := &http.Server{
		Addr:              socket,
		Handler:           api.Mux(),
		IdleTimeout:       3 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}

	err = srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
