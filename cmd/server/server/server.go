// Licensed under the Apache License, Version 2.0
// Details: https://raw.githubusercontent.com/square/quotaservice/master/LICENSE

package server

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/square/quotaservice"
	"github.com/square/quotaservice/buckets/memory"
	"github.com/square/quotaservice/config"
	"github.com/square/quotaservice/events"
	"github.com/square/quotaservice/logging"
	pb "github.com/square/quotaservice/protos/config"
	"github.com/square/quotaservice/rpc/grpc"
	"github.com/square/quotaservice/stats"
)

const (
	HTTPServer = "localhost:8080"
	gRPCServer = "localhost:10990"
)

func RunServer(cfg *pb.ServiceConfig) {
	server := quotaservice.New(memory.NewBucketFactory(),
		config.NewMemoryConfig(cfg),
		config.NewReaperConfig(),
		0,
		grpc.New(gRPCServer, events.NewNilProducer()))
	server.SetStatsListener(stats.NewMemoryStatsListener())
	if _, e := server.Start(); e != nil {
		panic(e)
	}

	// Serve Admin Console
	logging.Printf("Starting admin server on %v\n", HTTPServer)
	sm := http.NewServeMux()
	server.ServeAdminConsole(sm, "admin/public", true)
	go func() { _ = http.ListenAndServe(HTTPServer, sm) }()

	// Block until SIGTERM or SIGINT
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	var shutdown sync.WaitGroup
	shutdown.Add(1)

	go func() {
		<-sigs
		shutdown.Done()
	}()

	shutdown.Wait()
	_, _ = server.Stop()
}