package main

import (
	"flag"
	"log"
	"path/filepath"
	"runtime"

	"github.com/openagentsinc/v3/relay/internal/nip01"
	// "github.com/openagentsinc/v3/relay/internal/nip90"
	// "github.com/openagentsinc/v3/relay/internal/whisper"
	// "github.com/openagentsinc/v3/relay/internal/config"
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

	// TODO: Load configuration
	// config, err := config.Load()
	// if err != nil {
	//     log.Fatalf("Failed to load configuration: %v", err)
	// }

	// TODO: Initialize Whisper transcription service
	// whisperService, err := whisper.NewService(config.GroqAPIKey)
	// if err != nil {
	//     log.Fatalf("Failed to initialize Whisper service: %v", err)
	// }

	// TODO: Register Whisper transcription service with NIP-90
	// nip90Service := nip90.NewService()
	// nip90Service.RegisterProvider("whisper", whisperService)

	// TODO: Integrate NIP-90 service with the relay
	// relay.RegisterNIP90Service(nip90Service)

	// Start the WebSocket server
	log.Printf("Starting relay server on %s", *addr)
	err := relay.Start(*addr)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}