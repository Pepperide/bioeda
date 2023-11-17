package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/Pepperide/bioeda/operators/cmd/slurm-controller/proto"

	"google.golang.org/grpc"
)

const (
	ADDRESS       = "localhost:50051"
	GRPC_MAX_SIZE = 4194304
)

func main() {
	// Define flags
	path := flag.String("in", "", "path to the file")
	filename := flag.String("name", "defaultName", "Name of the output file")
	flag.Parse()

	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}

	defer conn.Close()

	c := pb.NewJobsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	/* ---- LOGIC ---- */

	messageSize := getMessageSize(*filename, *path)

	fileChunks := getFileChunks(*path, messageSize-len(*filename))

	_, err = c.Submit(ctx, &pb.JsonRequest{Filename: *filename, Json: fileChunks})
	if err != nil {
		log.Fatalf("Could not Submit script: %v", err)
	}
	log.Printf("Success. Message size is %v", messageSize)
}

func getMessageSize(filename string, path string) (size int) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening the file : %v", err)
	}

	stats, err := file.Stat()
	file.Close()

	s := int(stats.Size()) + len(filename)
	if s > GRPC_MAX_SIZE {
		log.Fatalf("Error: the size message exceed the maximum size (%v vs %v)", s, GRPC_MAX_SIZE)
	}

	return s
}

func getFileChunks(path string, size int) (buffer []byte) {

	buff := make([]byte, size)

	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("Error opening the file : %v", err)
	}

	_, err = file.ReadAt(buff, 0)

	if err != nil && err != io.EOF {
		log.Fatalf("Error reading file : %v", err)
	}

	file.Close()
	return buff
}
