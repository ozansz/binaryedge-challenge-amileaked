package main

import (
	"context"
	fmt "fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func runHTTPReverseProxy(httpEndpoint string, endpoint string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	fmt.Printf("Started Reverse Proxy Server to %s\n", endpoint)
	err := RegisterLeakServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	fmt.Printf("Started HTTP Server on %s\n", httpEndpoint)
	err = http.ListenAndServe(httpEndpoint, mux)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
}

func main() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChannel
		os.Exit(0)
	}()

	runHTTPReverseProxy(":8081", "localhost:50051")
}
