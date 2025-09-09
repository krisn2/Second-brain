package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krisn2/second-brain/db"
	"github.com/krisn2/second-brain/models"
)

func AddContent(c *gin.Context) {

	var body struct {
		Title   string   `json:"title"`
		Type    string   `json:"type"`
		Link    string   `json:"link"`
		Tags    []string `json:"tags"`
		Content string   `json:"content"`
	}

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid Input"})
		return
	}

	userID, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uuidUser, _ := uuid.Parse(userID.(string))

	var tags []models.Tag
	for _, t := range body.Tags {
		tag := models.Tag{Title: t}
		db.DB.FirstOrCreate(&tag, models.Tag{Title: t})
		tags = append(tags, tag)
	}

	content := models.Content{
		Link:    body.Link,
		Type:    body.Type,
		Title:   body.Title,
		UserId:  uuidUser,
		Content: body.Content,
		Tags:    tags,
	}

	if err := db.DB.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": "Brain created successfully",
		},
	)
}

func GetContent(c *gin.Context) {

	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
	}

	var content models.Content

	if err := db.DB.Where("title LIKE ?", "%"+query+"%").First(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	url := "https://api.groq.com/openai/v1/chat/completions"
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY environment variable not set")
	}

	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type ReqBody struct {
		Messages []Message `json:"messages"`
		Model    string    `json:"model"`
	}

	type Choice struct {
		Message Message `json:"message"`
	}

	type ResBody struct {
		Choices []Choice `json:"choices"`
	}

	requestBody := ReqBody{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Explan me this and give real world example of the content i provided and also give me a short summary of the content" + query,
			},
		},
		Model: "llama-3.3-70b-versatile",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to AI API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		c.JSON(resp.StatusCode, gin.H{"error": "AI API request failed", "details": errorBody})
		return
	}

	var resBody ResBody
	if err := json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
		return
	}

	aiResponseContent := ""
	if len(resBody.Choices) > 0 {
		aiResponseContent = resBody.Choices[0].Message.Content
	}

	c.JSON(http.StatusOK, gin.H{
		"ai_response":      aiResponseContent,
		"original_content": content,
	})
}
