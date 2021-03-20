package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"key_value/server/pb"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedKeyValueServiceServer
}

var storage = make(map[string]string)

func (s *server) Get(_ context.Context, in *pb.Key) (*pb.Value, error) {
	return &pb.Value{Value: storage[in.Key], Defined: storage[in.Key] != ""}, nil
}

func (s *server) Put(_ context.Context, in *pb.KeyValue) (*pb.Empty, error) {
	storage[in.Key] = in.Value

	return &pb.Empty{}, nil
}

func (s *server) GetAllKeys(_ context.Context, _ *pb.Empty) (*pb.StoredKeys, error) {
	var keys []string

	for k := range storage {
		keys = append(keys, k)
	}

	return &pb.StoredKeys{Keys: keys}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to listen - %s\n", err)
		os.Exit(1)
	}

	s := grpc.NewServer()

	fmt.Printf("Server listening on localhost%s\n", port)

	pb.RegisterKeyValueServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to serve - %s\n", err)
		os.Exit(1)
	}
}
