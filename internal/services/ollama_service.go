package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"vet-tails/ai/internal/models"
)

type LlavaService struct {
	baseURL string
	model   string
}

type LlavaRequest struct {
	Model       string   `json:"model"`
	Prompt      string   `json:"prompt"`
	Stream      bool     `json:"stream"`
	Temperature float32  `json:"temperature"`
	Images      []string `json:"images"`
}

type LlavaResponse struct {
	Response string `json:"response"`
}

func NewLlavaService(baseURL string) *LlavaService {
	return &LlavaService{
		baseURL: baseURL,
		model:   "llava", // or any medical-focused model available in Ollama
	}
}

func (s *LlavaService) DetectBreed(image string) (*models.BreedDetection, error) {
	// Prepare Ollama request
	prompt := `Analyze this pet image as a professional veterinarian:
    1. What breed do you see? Be specific.
    2. What visual characteristics support this identification?
    3. Rate your confidence level (0-100%).
    4. List any possible alternative breeds if unsure.
    
    Format your response as:
    Primary Breed: [breed name]
    Confidence: [percentage]
    Key Features: [list key identifying features]
    Alternative Breeds: [if applicable]`

	reqBody := LlavaRequest{
		Model:       "llava",
		Prompt:      prompt,
		Temperature: 0.7,
		Images:      []string{image},
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/generate", s.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("error calling Ollama API: %v", err)
	}
	defer resp.Body.Close()

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Parse the JSON response into Note structure
	var note models.BreedDetection
	if err := json.Unmarshal([]byte(result.Response), &note); err != nil {
		// Log the response for debugging
		fmt.Printf("Raw response: %s\n", result.Response)
		return nil, fmt.Errorf("error parsing SOAP note: %v", err)
	}

	return &note, nil
}
