package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Pepperide/bioeda/operators/cmd/slurm-controller-cluster/grpc/proto"
	logic "github.com/Pepperide/bioeda/operators/pkg/grpc/slurm-controller-cluster"
	"google.golang.org/grpc"
)

type JobsServer struct {
	pb.UnimplementedJobsServer
}

const (
	// Port for gRPC server to listen to
	PORT = ":50051"
)

func (s *JobsServer) Submit(ctx context.Context, req *pb.ScriptRequest) (*pb.SubmitResponse, error) {
	/* --- LOGIC --- */
	ok := logic.SubmitImpl(req.Filename, req.ScriptChunks)

	response := &pb.SubmitResponse{
		Ok: ok,
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
