package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/Pepperide/bioeda/operators/cmd/slurm-controller/proto"

	"google.golang.org/grpc"
)

const (
	ADDRESS       = "localhost:50051"
	GRPC_MAX_SIZE = 4000000
)

func main() {
	// Define flags
	path := flag.String("file", "", "path to the file")
	flag.Parse()

	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}

	defer conn.Close()

	c := pb.NewJobsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	// Logic
	file, err := os.Open(*path)

	if err != nil {
		log.Fatalf("Error opening the file : %v", err)
	}

	buff := make([]byte, GRPC_MAX_SIZE)

	n_bytes, err := file.ReadAt(buff, 0)
	fmt.Printf("Bytes read: %v\n", n_bytes)

	if err != nil && err != io.EOF {
		log.Fatalf("Error reading file : %v", err)
	}

	_, err = c.Submit(ctx, &pb.JsonRequest{Filename: "funzionaaaaa.json", Json: buff})
	if err != nil {
		log.Fatalf("Could not Submit script: %v", err)
	}

}
