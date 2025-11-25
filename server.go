package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"
)

// Message represents one chat entry
type Message struct {
	User   string
	Body   string
	SentAt time.Time
}

// ChatCore manages messages and chat history
type ChatCore struct {
	lock    sync.Mutex
	records []Message
}

// SendMessage adds a message to chat history and returns all messages
func (c *ChatCore) SendMessage(msg Message, reply *[]Message) error {
	if msg.User == "" {
		return errors.New("username cannot be empty")
	}
	if msg.Body == "" {
		return errors.New("message body cannot be empty")
	}
	if msg.SentAt.IsZero() {
		msg.SentAt = time.Now()
	}

	c.lock.Lock()
	c.records = append(c.records, msg)
	historyCopy := make([]Message, len(c.records))
	copy(historyCopy, c.records)
	c.lock.Unlock()

	// Print message to the server console
	timeStr := msg.SentAt.Format("15:04:05")
	if msg.User == "System" {
		// Cyan for system updates
		fmt.Printf("\033[36m[%s] %s\033[0m\n", timeStr, msg.Body)
	} else {
		// Green for usernames, white for messages
		fmt.Printf("\033[32m[%s] %s:\033[0m %s\n", timeStr, msg.User, msg.Body)
	}

	*reply = historyCopy
	return nil
}

// GetHistory returns the entire stored chat log
func (c *ChatCore) GetHistory(_ struct{}, reply *[]Message) error {
	c.lock.Lock()
	historyCopy := make([]Message, len(c.records))
	copy(historyCopy, c.records)
	c.lock.Unlock()

	*reply = historyCopy
	return nil
}

func main() {
	port := os.Getenv("CHAT_PORT")
	if port == "" {
		port = "1234"
	}

	server := new(ChatCore)
	if err := rpc.Register(server); err != nil {
		log.Fatalf("RPC registration failed: %v", err)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	fmt.Printf("\n\033[32mChat server running on port %s\033[0m\n", port)
	fmt.Println("Commands: 'quit' to stop | 'reset' to clear chat")

	// Command handler for admin console
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			cmd := strings.ToLower(strings.TrimSpace(scanner.Text()))
			switch cmd {
			case "quit":
				fmt.Println("Server shutting down...")
				listener.Close()
				return
			case "reset":
				server.lock.Lock()
				server.records = nil
				server.lock.Unlock()
				fmt.Println("Chat log cleared.")
			default:
				fmt.Println("Available commands: quit | reset")
			}
		}
	}()

	// Accept client connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			log.Printf("Connection error: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
