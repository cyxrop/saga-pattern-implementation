package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/cache"
	appGrpc "gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/grpc/invoice"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/mw"
	repository "gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/repository/invoice"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/service"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/db"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/queue/consumer"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/pkg/api"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/pkg/metrics"
	"google.golang.org/grpc"
)

const (
	EnvDBConn      = "DB_CONN"
	EnvBrokers     = "BROKERS"
	EnvMetricsHost = "METRICS_HOST"
	EnvGrpcHost    = "GRPC_HOST"
	EnvCacheHosts  = "CACHE_HOSTS"
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

	grpcHost := os.Getenv(EnvGrpcHost)
	if grpcHost == "" {
		log.Fatalf("%s is empty", EnvGrpcHost)
	}

	cacheHosts := os.Getenv(EnvCacheHosts)
	if grpcHost == "" {
		log.Fatalf("%s is empty", EnvCacheHosts)
	}

	ctx, cancel := context.WithCancel(context.Background())
	dbConn, err := db.New(ctx, dbConnStr)
	if err != nil {
		log.Fatalf("db conn: %s", err)
	}
	defer dbConn.Close()

	c := cache.NewCache(strings.Split(cacheHosts, ";"))
	invoicesRepo := repository.New(dbConn)
	invoiceService := service.New(invoicesRepo, c)

	invoiceServer := appGrpc.NewInvoiceServiceServer(invoiceService)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(mw.LogInterceptor))
	api.RegisterInvoiceServiceServer(grpcServer, invoiceServer)

	lis, err := net.Listen("tcp", grpcHost)
	if err != nil {
		log.Printf("listen: %v\n", err)
		return
	}

	if err = consumer.Consume(ctx, strings.Split(brokers, ";"), invoiceService); err != nil {
		fmt.Printf("run consume error: %s", err)
		return
	}

	errCh := make(chan error)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	runMetricsSrv(errCh, metricsHost)

	go func() {
		defer close(errCh)
		log.Printf("Run GRPC server in %q...\n", grpcHost)
		if err = grpcServer.Serve(lis); err != nil {
			errCh <- fmt.Errorf("serve grpc server failed: %w", err)
		}
	}()

	select {
	case <-signals:
		log.Println("Graceful shutdown on syscall")
		cancel()
	case err = <-errCh:
		log.Printf("Graceful shutdown on err: %s", err)
		cancel()
	}
}

func runMetricsSrv(errCh chan<- error, host string) {
	expvar.Publish("Goroutines", &metrics.Goroutines{})

	srv := &http.Server{Addr: host}
	go func() {
		defer close(errCh)

		if err := srv.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("serve metrics server failed: %w", err)
		}
	}()
}
