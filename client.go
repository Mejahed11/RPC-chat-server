package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
)

// Message defines a single message in the chat
type Message struct {
	User   string
	Body   string
	SentAt time.Time
}

// renderChat prints all messages with color-coded formatting
func renderChat(logs []Message, currentUser string) {
	fmt.Println("\n--- Conversation ---")
	for _, entry := range logs {
		timestamp := entry.SentAt.Format("15:04:05")

		switch {
		case entry.User == "System":
			// 🩵 Cyan for system messages
			fmt.Printf("\033[36m[%s] %s\033[0m\n", timestamp, entry.Body)
		case entry.User == currentUser:
			// 🟡 Yellow for your messages
			fmt.Printf("\033[33m[%s] You: %s\033[0m\n", timestamp, entry.Body)
		default:
			// 🔴 Red for others
			fmt.Printf("\033[31m[%s] %s: %s\033[0m\n", timestamp, entry.User, entry.Body)
		}
	}
	fmt.Println("---------------------")
}

func main() {
	server := os.Getenv("CHAT_SERVER")
	if server == "" {
		server = "localhost:1234"
	}

	client, err := rpc.Dial("tcp", server)
	if err != nil {
		log.Fatalf("❌ Could not connect to %s: %v", server, err)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter nickname: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)
	if user == "" {
		user = "Guest"
	}

	fmt.Printf("\n\033[36mConnected as %s. Welcome!\033[0m\n", user)

	// Notify others
	joinMsg := Message{
		User:   "System",
		Body:   fmt.Sprintf("%s joined the chatroom.", user),
		SentAt: time.Now(),
	}

	var history []Message
	client.Call("ChatCore.SendMessage", joinMsg, &history)
	renderChat(history, user)

	for {
		fmt.Print("Type message (or 'quit' to exit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "quit") {
			leaveMsg := Message{
				User:   "System",
				Body:   fmt.Sprintf("%s has disconnected.", user),
				SentAt: time.Now(),
			}
			client.Call("ChatCore.SendMessage", leaveMsg, &history)
			fmt.Println("\033[36mYou have left the chat.\033[0m")
			return
		}

		if input == "" {
			continue
		}

		newEntry := Message{
			User:   user,
			Body:   input,
			SentAt: time.Now(),
		}

		history = nil
		if err := client.Call("ChatCore.SendMessage", newEntry, &history); err != nil {
			log.Printf("Server issue: %v", err)
			fmt.Println("⚠️ Connection lost. Try again later.")
			return
		}

		renderChat(history, user)
	}
}
