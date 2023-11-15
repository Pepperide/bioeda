package main

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"

	pb "github.com/Pepperide/bioeda/operators/cmd/slurm-controller/proto"
	"google.golang.org/grpc"
)

const (
	// Port for gRPC server to listen to
	PORT = ":50051"
)

type JobsServer struct {
	pb.UnimplementedJobsServer
}

func (s *JobsServer) Submit(ctx context.Context, json_req *pb.JsonRequest) (*pb.SubmitResponse, error) {
	// Business logic
	log.Println("Submitting...")
	filename := filepath.Join("/home/control", json_req.Filename)
	err := os.WriteFile(filename, json_req.Json, 0600)
	if err != nil {
		log.Fatalf("Could not Write file: %v", err)
	}
	response := &pb.SubmitResponse{
		Ok: err == nil,
	}

	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterJobsServer(s, &JobsServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
