# Relay Implementation Plan

This document outlines the step-by-step plan for setting up our relay folder to implement NIP-01 and NIP-90 from scratch in Golang, with the first NIP-90 service being audio transcription using Whisper via the Groq API.

## 1. Project Setup

1. Initialize a new Go module in the `relay` folder:
   ```
   cd relay
   go mod init github.com/openagentsinc/v3/relay
   ```

2. Create the basic directory structure:
   ```
   mkdir -p cmd/relay
   mkdir -p internal/nostr
   mkdir -p internal/nip01
   mkdir -p internal/nip90
   mkdir -p internal/whisper
   ```

## 2. Implement NIP-01 (Basic Protocol Flow)

1. Create `internal/nostr/event.go`:
   - Define the `Event` struct
   - Implement event serialization and deserialization

2. Create `internal/nostr/filter.go`:
   - Define the `Filter` struct
   - Implement filter matching logic

3. Create `internal/nip01/messages.go`:
   - Define structs for different message types (EVENT, REQ, CLOSE, etc.)
   - Implement message parsing and creation functions

4. Create `internal/nip01/subscription.go`:
   - Implement subscription management

5. Create `internal/nip01/relay.go`:
   - Implement the main relay logic
   - Handle WebSocket connections
   - Process incoming messages
   - Manage subscriptions
   - Broadcast events to subscribers

## 3. Implement NIP-90 (Data Vending Machines)

1. Create `internal/nip90/job.go`:
   - Define structs for Job Request, Job Result, and Job Feedback
   - Implement job creation, processing, and feedback handling

2. Create `internal/nip90/service.go`:
   - Define the interface for NIP-90 services
   - Implement service registration and discovery

3. Create `internal/nip90/whisper_service.go`:
   - Implement the Whisper transcription service using the Groq API

4. Update `internal/nip01/relay.go`:
   - Integrate NIP-90 job handling into the main relay logic

## 4. Implement Whisper Transcription Service

1. Create `internal/whisper/client.go`:
   - Implement a client for the Groq API
   - Handle authentication and API requests

2. Create `internal/whisper/transcribe.go`:
   - Implement the audio transcription logic using the Groq API client

## 5. Main Application

1. Create `cmd/relay/main.go`:
   - Set up the main application entry point
   - Initialize the relay
   - Register the Whisper transcription service
   - Start the WebSocket server

## 6. Configuration and Environment

1. Create `internal/config/config.go`:
   - Implement configuration loading from environment variables or config file
   - Include settings for Groq API key, relay URL, etc.

## 7. Testing

1. Create test files for each package:
   - `internal/nostr/event_test.go`
   - `internal/nip01/relay_test.go`
   - `internal/nip90/job_test.go`
   - `internal/whisper/transcribe_test.go`

2. Implement unit tests for core functionality

3. Create integration tests in `test/integration_test.go`

## 8. Documentation

1. Update `README.md` with:
   - Project overview
   - Setup instructions
   - Usage examples
   - API documentation

2. Add inline documentation to all exported functions and types

## 9. Deployment

1. Create a `Dockerfile` for containerization

2. Set up CI/CD pipeline (e.g., GitHub Actions) for automated testing and deployment

## Next Steps

- Implement additional NIP-90 services
- Enhance error handling and logging
- Implement performance optimizations
- Add monitoring and metrics
- Consider splitting the relay and NIP-90 service provider into separate components

Remember to commit your changes regularly and follow Go best practices and conventions throughout the implementation.