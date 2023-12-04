package slurmcontrollercluster

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/Pepperide/bioeda/operators/cmd/slurm-controller-cluster/grpc/proto"
)

const (
	GRPC_MAX_SIZE = 4194304
)

func SendMessage(c pb.JobsClient, ctx context.Context, filepath string) error {

	messageSize, chunksSize, err := getMessageSize(filepath)
	if err != nil {
		return err
	}

	fileChunks, err := getFileChunks(filepath, chunksSize)
	if err != nil {
		return err
	}

	_, err = c.Submit(ctx, &pb.ScriptRequest{Filename: filepath, ScriptChunks: fileChunks})
	if err != nil {
		return err
	}

	log.Printf("Success. Message size is %v", messageSize)

	return nil
}

func getMessageSize(path string) (int, int, error) {

	file, err := os.Open(path)
	if err != nil {
		return -1, -1, err
	}

	stats, err := file.Stat()
	if err != nil {
		return -1, -1, err
	}

	s := int(stats.Size()) + len(stats.Name())
	if s > GRPC_MAX_SIZE {
		return -1, -1, fmt.Errorf("Error: the size message exceed the maximum size (%v vs %v)", s, GRPC_MAX_SIZE)
	}

	file.Close()
	return s, len(stats.Name()), nil
}

func getFileChunks(path string, size int) ([]byte, error) {

	buff := make([]byte, size)

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	_, err = file.ReadAt(buff, 0)

	if err != nil && err != io.EOF {
		return nil, err
	}

	file.Close()
	return buff, nil
}
