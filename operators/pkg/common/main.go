package common

import (
	"flag"
	"log"
)

func IsFlagPassed(name string) bool {
	log.Printf("Check flags")
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
