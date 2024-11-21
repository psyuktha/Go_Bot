package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

func printResponse(resp *genai.GenerateContentResponse) string {
	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		result += fmt.Sprintf("%v\n", part)
	}
	return result
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

	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/chat", func(c *gin.Context) {
		var json struct {
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := sendMessage(ctx, session, json.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		responseText := printResponse(resp)
		c.JSON(http.StatusOK, gin.H{"response": responseText})
	})

	r.Run(":8080")
}
