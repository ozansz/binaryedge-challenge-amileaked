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
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		log.Fatal(err)
	}

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

	handler := &LeakServiceServerHandler{
		DBConnURI:    "mongodb+srv://dbUser:uJKmMse-U3Tgk5K@cluster0-m6gsh.mongodb.net/?ssl=true",
		DatabaseName: "ail",
	}

	err = handler.DBConnect()

	if err != nil {
		log.Fatal(err)
	}

	RegisterLeakServiceServer(srv, handler)

	reflection.Register(srv)

	fmt.Printf("Started GRPC Server on %s\n", fmt.Sprintf("%s:%d", host, port))
	err = srv.Serve(listener)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChannel
		os.Exit(0)
	}()

	runGRPCServer("localhost", 50051)
}
