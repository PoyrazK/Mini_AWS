package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/poyrazk/thecloud/internal/storage/node"
	pb "github.com/poyrazk/thecloud/internal/storage/protocol"
	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "9101", "Port to listen on")
	dataDir := flag.String("data-dir", "./data/storage-node", "Directory to store data")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting storage node", "port", *port, "dataDir", *dataDir)

	// 1. Init Store
	store, err := node.NewLocalStore(*dataDir)
	if err != nil {
		logger.Error("failed to init store", "error", err)
		os.Exit(1)
	}

	// 2. Init RPC Server
	rpcServer := node.NewRPCServer(store)
	grpcServer := grpc.NewServer()
	pb.RegisterStorageNodeServer(grpcServer, rpcServer)

	// 3. Listen
	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		logger.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	// 4. Handle Shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		logger.Info("shutting down")
		grpcServer.GracefulStop()
	}()

	logger.Info("storage node ready")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
