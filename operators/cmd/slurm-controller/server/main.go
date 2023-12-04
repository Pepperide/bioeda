/* --- SLURM CONTROLLER --- */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	pb "github.com/Pepperide/bioeda/operators/cmd/slurm-controller/proto"
	"google.golang.org/grpc"
)

const (
	// Port for gRPC server to listen to
	PORT       = ":50051"
	TMP_FOLDER = "/tmp"
)

type JobsServer struct {
	pb.UnimplementedJobsServer
}

type SBatchOpt struct {
	Output string
	Cpus   int
	Nodes  int
}

type Script struct {
	Name    string
	Command string
}

type Job struct {
	UserID  string
	JobName string
	Sbatch  SBatchOpt
	Scripts []Script
}

func (s *JobsServer) Submit(ctx context.Context, json_req *pb.JsonRequest) (*pb.SubmitResponse, error) {

	/* --- LOGIC --- */

	// Parse the json script sent by the client
	job := jsonParser(json_req.Json)

	// Create the script and save it
	ok := makeScript(job)

	response := &pb.SubmitResponse{
		Ok: ok,
	}

	return response, nil
}

func jsonParser(chunks []byte) Job {
	jobString := string(chunks[:])
	fmt.Println(jobString)

	var job Job
	err := json.Unmarshal([]byte(jobString), &job)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	return job
}

func makeScript(job Job) bool {
	sbatchOpt := job.Sbatch

	// Save the script for testing
	path := filepath.Join("/home", "", job.UserID)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}

	file, err := os.Create(filepath.Join(path, job.JobName+".sh"))
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}

	// File header
	str := "#!/bin/bash"

	// Options
	if sbatchOpt.Cpus != 0 {
		str += "\n#SBATCH --cpus-per-task=" + fmt.Sprintf("%d", sbatchOpt.Cpus)
	}
	if sbatchOpt.Nodes != 0 {
		str += "\n#SBATCH --nodes=" + fmt.Sprintf("%d", sbatchOpt.Nodes)
	}
	if sbatchOpt.Output != "" {
		str += "\n#SBATCH --output=" + sbatchOpt.Output
	}

	// Script
	for _, value := range job.Scripts {
		str += "\n" + value.Command + value.Name
	}

	// Write file
	_, err = file.Write([]byte(str))
	if err != nil {
		log.Printf("Could not Write file: %v", err)
		return false
	}

	return true
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
