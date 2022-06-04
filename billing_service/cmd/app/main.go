package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	repository "gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/repository/invoice"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/service"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/db"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/queue/consumer"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/pkg/metrics"
)

const (
	EnvDBConn      = "DB_CONN"
	EnvBrokers     = "BROKERS"
	EnvMetricsHost = "METRICS_HOST"
)

func main() {
	log.Println("Start billing service...")

	dbConnStr := os.Getenv(EnvDBConn)
	if dbConnStr == "" {
		log.Fatalf("%s is empty", EnvDBConn)
	}

	brokers := os.Getenv(EnvBrokers)
	if brokers == "" {
		log.Fatalf("%s is empty", EnvBrokers)
	}

	metricsHost := os.Getenv(EnvMetricsHost)
	if metricsHost == "" {
		log.Fatalf("%s is empty", EnvMetricsHost)
	}

	ctx, cancel := context.WithCancel(context.Background())
	dbConn, err := db.New(ctx, dbConnStr)
	if err != nil {
		log.Fatalf("db conn: %s", err)
	}
	defer dbConn.Close()

	ir := repository.New(dbConn)
	invoiceService := service.New(ir)

	if err = consumer.Consume(ctx, strings.Split(brokers, ";"), invoiceService); err != nil {
		fmt.Printf("run consume error: %s", err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-signals:
		log.Println("Graceful shutdown on syscall")
		cancel()
	case err = <-runMetricsSrv(metricsHost):
		log.Printf("Graceful shutdown on mitrics server err: %s", err)
		cancel()
	}
}

func runMetricsSrv(host string) <-chan error {
	expvar.Publish("Goroutines", &metrics.Goroutines{})

	errs := make(chan error)
	srv := &http.Server{Addr: host}

	go func() {
		defer close(errs)

		if err := srv.ListenAndServe(); err != nil {
			errs <- err
		}
	}()

	return errs
}
