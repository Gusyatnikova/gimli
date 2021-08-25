package  main

import (
	"gimli/Internal/gRPC/domain"
	"gimli/Internal/gRPC/impl"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	log.Println("gRPC In Action...")
	s := grpc.NewServer()
	shortenerServ := impl.ShortenerServiceImpl{}
	domain.RegisterShortenerServer(s, shortenerServ)

	l, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatal(err)
	}
	shortenerServ.Begin()
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
	shortenerServ.End()
}