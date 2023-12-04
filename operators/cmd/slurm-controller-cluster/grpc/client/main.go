package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/Pepperide/bioeda/operators/cmd/slurm-controller-cluster/grpc/proto"
	logic "github.com/Pepperide/bioeda/operators/pkg/grpc/slurm-controller-cluster"

	"google.golang.org/grpc"
)

const (
	ADDRESS = "localhost:50051"
)

func main() {
	// Define flags
	filepath := flag.String("filepath", "", "path to the file")
	// filename := flag.String("name", "defaultName", "Name of the output file")
	flag.Parse()

	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}

	defer conn.Close()

	c := pb.NewJobsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	err = logic.SendMessage(c, ctx, *filepath)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}
