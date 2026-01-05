# TCP Async Image Processor

This project implements a raw TCP server designed to handle asynchronous image processing tasks. It demonstrates concurrency patterns in Go using worker pools, buffered channels, and mutexes to manage resources efficiently under load without relying on external frameworks.

## Project Structure

- **cmd/**: Entry points for the server and client applications.
- **server/**: Core server logic including connection handling, worker pools, and job queues.
- **client/**: Client implementation for reading images and sending them over TCP.
- **spam.sh**: A shell script to simulate high load by spawning multiple concurrent clients.

## Building the Project

To build the executable binaries for both the server and the client, run the following commands from the root directory:

```bash
go build -o server_bin cmd/server/main.go
go build -o client_bin cmd/client/main.go
```

## Running the Application

### 1. Start the Server

First, start the TCP server. It will listen on port 8082.

```bash
./server_bin
```

The server initializes a fixed number of workers and waits for incoming connections. It uses a semaphore pattern to limit the maximum number of concurrent connections.

### 2. Run a Single Client

To run a single client that reads images from the source directory and sends them to the server:

```bash
./client_bin
```

### 3. Load Testing

To simulate concurrent load and test the server's backpressure and concurrency limits, you can use the provided script. This script launches 70 client instances simultaneously.

```bash
./spam.sh
```

## Technical Details

The server does not spawn a new goroutine for every request indefinitely. Instead, it uses:
- **Worker Pool**: A fixed number of background workers that process jobs from a queue.
- **Connection Limiting**: A buffered channel acts as a semaphore to hard-limit active connections, preventing resource exhaustion.
- **Thread Safety**: Mutexes are used to safely track batch progress and metrics across multiple workers.
