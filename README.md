# WebSocket Server for Text Streaming

This repository contains a Go-based WebSocket server that dynamically switches between multiple inference providers based on response time and error rates. It also includes a Python script to test the WebSocket server by simulating multiple client connections.

## Features

- **Dynamic Provider Switching**: The server dynamically switches between inference providers if a provider exceeds a response time threshold or has a high error rate.
- **Error Handling**: Handles WebSocket connection errors gracefully, ensuring stable communication with clients.
- **Multiple Client Support**: The Python test client can simulate multiple clients connecting to the WebSocket server simultaneously, sending messages, and receiving responses.

## Contents

- **Go WebSocket Server (`websocketserver.go`)**: The main WebSocket server code.
- **Inference Logic (`inference.go`)**: Contains the logic for handling multiple inference providers.
- **Python Test Client (`ConnectionTest.py`)**: A Python script that connects multiple clients to the WebSocket server to test its functionality.

## Prerequisites

- **Go**: Make sure you have Go installed. [Download Go](https://golang.org/dl/)
- **Python**: Ensure Python 3.x is installed. [Download Python](https://www.python.org/downloads/)

## Installation

### Clone the Repository

```bash
git clone https://github.com/yourusername/websocket-server.git
cd TextStreamingEndpoint
```

### Run the Server

```bash
go run main.go
```

### Test the endpoint

- Open another Terminal
```bash
Python3 ConnectionTest.py
```

