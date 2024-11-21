package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Session struct {
	model *genai.GenerativeModel
}

func (s *Session) SendMessage(ctx context.Context, input genai.Text) (*genai.GenerateContentResponse, error) {
	return s.model.GenerateContent(ctx, input)
}

func sendMessage(ctx context.Context, session *Session, input string) (*genai.GenerateContentResponse, error) {
	resp, err := session.SendMessage(ctx, genai.Text(input))
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return nil, err
	}
	return resp, nil
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
}

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	session := &Session{model: model}

	// Example call to model.GenerateContent
	resp, err := model.GenerateContent(ctx, genai.Text("Write a story about a magic backpack."))
	if err != nil {
		log.Fatal(err)
	}
	printResponse(resp)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Start chatting with the model (type 'exit' to quit):")

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "exit" {
			break
		}

		resp, err := sendMessage(ctx, session, input)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			continue
		}

		fmt.Print("Model: ")
		printResponse(resp)
	}
}
