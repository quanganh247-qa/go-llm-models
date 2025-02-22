package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"vet-tails/ai/internal/models"
)

type SOAPService struct {
	ollamaURL string
}

type OllamaRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	Stream      bool    `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func NewSOAPService(ollamaURL string) *SOAPService {
	return &SOAPService{
		ollamaURL: ollamaURL,
	}
}
func (s *SOAPService) GenerateSOAPNote(transcribedText string) (*models.Note, error) {
	prompt := fmt.Sprintf(`As a veterinary AI assistant, analyze the following consultation transcript and generate a SOAP note:

    Transcript:
    %s

    Generate a structured SOAP note with the following sections:
    - Subjective (patient info, chief complaint, duration, history, symptoms)
    - Objective (vital signs, examination findings)
    - Assessment (primary diagnosis, differential diagnoses)
    - Plan (immediate treatment, medications, follow-up, client education)

    Format the response in a valid JSON structure matching this example:
    {
        "subjective": {
            "patient_info": "species, age, sex",
            "chief_complaint": "main issue",
            "duration": "time period",
            "history": "relevant history",
            "symptoms": ["symptom1", "symptom2"]
        },
        "objective": {
            "vital_signs": {
                "temperature": "value",
                "general_condition": "description"
				"respiratory_rate": "value",
				"weight": "value",
            },
            "examination_findings": ["finding1", "finding2"]
        },
        "assessment": {
            "primary_diagnosis": "main diagnosis",
            "differentials": ["differential1", "differential2"]
        },
        "plan": {
            "immediate_treatment": ["treatment1", "treatment2"],
            "medications": ["medication1", "medication2"],
            "follow_up": "follow up plan",
            "client_education": ["education1", "education2"]
        }
    }`, transcribedText)

	reqBody := OllamaRequest{
		Model:       "mistral",
		Prompt:      prompt,
		Temperature: 0.7,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/generate", s.ollamaURL),
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
	var note models.Note
	if err := json.Unmarshal([]byte(result.Response), &note); err != nil {
		// Log the response for debugging
		fmt.Printf("Raw response: %s\n", result.Response)
		return nil, fmt.Errorf("error parsing SOAP note: %v", err)
	}

	return &note, nil
}

// ... existing code ...

func (s *SOAPService) GeneratePatientSummary(patientHistory string) (*models.PatientSummary, error) {
	prompt := fmt.Sprintf(`As a veterinary AI assistant, create a concise patient summary from the following medical history:

    Medical History:
    %s

    Generate a brief, structured summary that a veterinarian can quickly review before entering the exam room. Include:
    - Key medical conditions
    - Recent visits and their purposes
    - Current medications
    - Important alerts (allergies, behavioral notes)
    - Preventive care status

    Format the response in a valid JSON structure matching this example:
    {
        "key_conditions": ["condition1", "condition2"],
        "recent_visits": [
            {
                "date": "YYYY-MM-DD",
                "reason": "visit reason"
            }
        ],
        "current_medications": [
            {
                "name": "medication name",
                "dosage": "dosage info",
                "frequency": "frequency info"
            }
        ],
        "alerts": ["alert1", "alert2"],
        "preventive_care": {
            "vaccinations": ["status1", "status2"],
            "next_due": ["due1", "due2"]
        }
    }`, patientHistory)

	reqBody := OllamaRequest{
		Model:       "mistral",
		Prompt:      prompt,
		Temperature: 0.7,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/generate", s.ollamaURL),
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

	var summary models.PatientSummary
	if err := json.Unmarshal([]byte(result.Response), &summary); err != nil {
		fmt.Printf("Raw response: %s\n", result.Response)
		return nil, fmt.Errorf("error parsing patient summary: %v", err)
	}

	return &summary, nil
}

func (s *SOAPService) GeneratePetActivityLog(activityDescription string) (*models.ActivityLog, error) {
	prompt := fmt.Sprintf(`As a veterinary AI assistant, analyze the following pet activity description and generate a structured activity log:

    Activity Description:
    %s

    Generate a detailed activity log with the following information:
    - Activity type (exercise, feeding, medication, grooming, behavior, etc.)
    - Duration and timing
    - Observations and notes
    - Any concerns or follow-up needed

    Format the response in a valid JSON structure matching this example:
    {
        "activity_type": "type of activity",
        "timestamp": "YYYY-MM-DD HH:MM",
        "duration": "duration in minutes",
        "details": {
            "description": "detailed description of activity",
            "intensity_level": "low/medium/high",
            "location": "where activity took place"
        },
        "observations": ["observation1", "observation2"],
        "concerns": ["concern1", "concern2"],
        "follow_up_needed": false,
        "notes": "additional relevant information"
    }`, activityDescription)

	reqBody := OllamaRequest{
		Model:       "mistral",
		Prompt:      prompt,
		Temperature: 0.7,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/api/generate", s.ollamaURL),
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

	var activityLog models.ActivityLog
	if err := json.Unmarshal([]byte(result.Response), &activityLog); err != nil {
		fmt.Printf("Raw response: %s\n", result.Response)
		return nil, fmt.Errorf("error parsing activity log: %v", err)
	}

	return &activityLog, nil
}
