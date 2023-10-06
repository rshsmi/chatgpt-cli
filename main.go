package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"bufio"
)

// Define a struct to hold a single message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Define a struct to hold the API request payload
type Payload struct {
	Model   string    `json:"model"`
	Messages []Message `json:"messages"`
}

func main() {
	// Create a new HTTP client
	client := &http.Client{}

	// Set up a loop to continuously prompt the user for input
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You: ")
		scanner.Scan()  // Read user input from CLI
		userMessage := scanner.Text()

		// Exit loop if user types "exit"
		if userMessage == "exit" {
			break
		}

		// Prepare the API request payload
		payload := Payload{
			Model: "gpt-3.5-turbo",
			Messages: []Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: userMessage},
			},
		}

		// Marshal the payload to JSON
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshalling payload:", err)
			return
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the necessary headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer YOUR_API_KEY")

		// Send the request and get the response
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Read and print the raw API response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Println("API Response:", string(body))
	}
}
