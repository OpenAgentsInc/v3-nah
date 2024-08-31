# Relay

This is the relay implementation for the OpenAgents project, supporting NIP-01 and NIP-90 protocols.

## Development

### Prerequisites

- Go 1.16 or later

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/openagentsinc/v3.git
   cd v3/relay
   ```

2. Install dependencies:
   ```
   go mod download
   ```

### Running the Relay

To run the relay in development mode:

```
go run cmd/relay/main.go
```

By default, the relay will start on port 8080. You can specify a different port using the `-addr` flag:

```
go run cmd/relay/main.go -addr :9000
```

### Testing

To run the tests:

```
go test ./...
```

## Building

To build the relay:

```
go build -o relay cmd/relay/main.go
```

This will create an executable named `relay` in your current directory.

## Deployment

To run the built relay:

```
./relay
```

Or with a custom address:

```
./relay -addr :9000
```

## Configuration

(TODO: Add information about configuration options and how to set them)

## Contributing

(TODO: Add information about how to contribute to the project)

## License

(TODO: Add license information)