package slurmcontrollercluster

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	SCRIPT_FOLDER = "/script"
)

func SubmitImpl(filename string, chunks []byte) bool {
	ok := true

	path, err := createFile(filename, chunks)
	if err != nil {
		log.Printf("%v", err)
		ok = false
	}

	err = startJob(path)
	if err != nil {
		log.Printf("%v", err)
		ok = false
	} else {
		log.Printf("Job started")
	}

	return ok
}

func createFile(filename string, chunks []byte) (string, error) {
	path := filepath.Join(SCRIPT_FOLDER, filename)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}

	_, err = file.Write(chunks)
	if err != nil {
		return "", err
	}

	file.Close()
	return path, nil
}

func startJob(path string) error {
	cmd := exec.Command("sbatch " + path)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
