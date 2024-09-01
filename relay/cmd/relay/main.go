package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/openagentsinc/v3/relay/internal/nip01"
)

func init() {
	// Change the working directory to the project root
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Parse command-line flags
	addr := flag.String("addr", ":8080", "HTTP service address")
	flag.Parse()

	// Initialize the relay
	relay := nip01.NewRelay()

	// Start the WebSocket server
	log.Printf("Starting relay server on %s", *addr)
	err := relay.Start(*addr)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
