# Dockerized RPC Chatroom (Go)

This project is a **Go RPC chatroom** packaged inside a **Docker container**, allowing you to run the server locally, connect with a Go client, and publish the Docker image on **Docker Hub**.

> **Docker Hub Image:**  
> `https://hub.docker.com/r/mohamedmejahed/rpc-chat-server`

---

## Features

### ✅ RPC Server
- Stores all chat messages **in-memory**.
- Exposes RPC methods:
  - `AddMessage` — adds a new message  
  - `GetHistory` — returns the full chat history
- Runs inside a **Docker container** on port `1234`.

### ✅ RPC Client
- Connects to the server at `localhost:1234`
- Sends user messages
- Receives and prints the chat history
- Keeps running until manually terminated

---

## How to Run (Dockerized Server + Local Client)

### 1) Run the Server (Docker)

```bash
docker pull mohamedmejahed/rpc-chat-server:new
docker run --rm -p 1234:1234 mohamedmejahed/rpc-chat-server:new
```

Server listens on:

```
localhost:1234
```

---

### 2) Run the Client (Go)

Open another terminal in the project folder:

```bash
go run client.go
```

Example interaction:

```
Enter message: Hello
--- Chat History ---
You: Hello
--------------------
```

---

## Dockerfile

```dockerfile
FROM golang:1.22-alpine
WORKDIR /app
COPY server.go .
RUN go build -o server server.go
ENV CHAT_PORT=1234
EXPOSE 1234
CMD ["./server"]
```

---

## Image Build and Publish

### 1) Build the image
```bash
docker build -t mohamedmejahed/rpc-chat-server:new .
```

### 2) Login to Docker Hub
```bash
docker login
```

### 3) Push the image
```bash
docker push mohamedmejahed/rpc-chat-server:new
```
---

