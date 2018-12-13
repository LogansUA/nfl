package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/logansua/nfl_app/bucket"
	"github.com/logansua/nfl_app/db"
	"github.com/logansua/nfl_app/player"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	dbService := db.New("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PASSWORD"),
	))
	bucketService := bucket.New()
	playerService := player.New(dbService, bucketService)

	var handler http.Handler
	{
		handler = player.MakeHTTPHandler(playerService, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal)

		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		errs <- fmt.Errorf("%playerService", <-c)
	}()

	httpAddr := flag.String("http.addr", fmt.Sprintf(":%s", os.Getenv("APP_PORT")), "HTTP listen address")

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)

		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	logger.Log("exit", <-errs)
}
