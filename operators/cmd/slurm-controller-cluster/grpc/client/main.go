package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/Pepperide/bioeda/operators/cmd/slurm-controller-cluster/grpc/proto"
	common "github.com/Pepperide/bioeda/operators/pkg/common"
	logic "github.com/Pepperide/bioeda/operators/pkg/grpc/slurm-controller-cluster"

	"google.golang.org/grpc"
)

var domain = "localhost"
var address = domain + ":50051"

func main() {
	// Define flags
	var filepath string
	flag.StringVar(&filepath, "f", "", "path to the file")
	// filename := flag.String("name", "defaultName", "Name of the output file")
	flag.Parse()
	log.Printf("File: %v", filepath)
	found := common.IsFlagPassed("f")
	if !found {
		log.Printf("Error: the file path is not set")
		return
	}

	log.Printf("Try to connect %v", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Did not connect : %v", err)
		return
	}

	defer conn.Close()

	c := pb.NewJobsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	err = logic.SendMessage(c, ctx, filepath)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}
