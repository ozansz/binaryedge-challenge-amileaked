package main

import (
	fmt "fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runGRPCServer(host string, port int64) {
	// Open socket to listen on
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		log.Fatal(err)
	}

	// Create gRPC server with inceptor middleware
	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// grpc_zap.StreamServerInterceptor(zapLogger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// grpc_zap.UnaryServerInterceptor(zapLogger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	// Construct the service server
	handler := &LeakServiceServerHandler{
		DBConnURI:    "mongodb+srv://dbUser:passwd@localhost/?ssl=true",
		DatabaseName: "ail",
	}

	// Make the DB connection
	err = handler.DBConnect()

	if err != nil {
		log.Fatal(err)
	}

	// Register server handler to the server created
	RegisterLeakServiceServer(srv, handler)
	reflection.Register(srv)

	fmt.Printf("Started GRPC Server on %s\n", fmt.Sprintf("%s:%d", host, port))

	// Serve on the socket opened before
	err = srv.Serve(listener)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// Wait for keyboard interrupt to exit the program
	go func() {
		<-sigChannel
		os.Exit(0)
	}()

	// Let the fun begin!
	runGRPCServer("localhost", 50051)
}
