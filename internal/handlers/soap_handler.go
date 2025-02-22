package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"vet-tails/ai/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Input struct {
	PatientID  uint   `json:"patient_id"`
	Transcript string `json:"transcript"`
	Summary    string `json:"summary"`
}

type Output struct {
	SOAPNote string `json:"soap_note"`
}

type BreedOutput struct {
	Breed string `json:"breed"`
}

type Handler struct {
	DB           *gorm.DB
	LlavaService *services.LlavaService
	SoapService  *services.SOAPService
}

func (h *Handler) CreateSOAPNote(c *gin.Context) {

	var input Input

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.SoapService.GenerateSOAPNote(input.Transcript)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"soap_note": note,
	})

}

func (h *Handler) GeneratePatientSummary(c *gin.Context) {
	var input Input

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summary, err := h.SoapService.GeneratePatientSummary(input.Summary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}

func (h *Handler) DetectBreed(c *gin.Context) {

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Read the file content
	fileBytes, err := io.ReadAll(image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to base64 string
	base64String := base64.StdEncoding.EncodeToString(fileBytes)

	breed, err := h.LlavaService.DetectBreed(base64String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"breed": breed,
	})
}

func (h *Handler) GeneratePetActivityLog(c *gin.Context) {
	var input Input

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activityLog, err := h.SoapService.GeneratePetActivityLog(input.Transcript)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activity_log": activityLog,
	})
}

type UploadPDFInput struct {
	Collection string `json:"collection"`
}

func (h *Handler) UploadPDFHandler(c *gin.Context) {

	file, err := c.FormFile("pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No PDF file provided"})
		return
	}

	// Save the uploaded file temporarily
	tempPath := fmt.Sprintf("/tmp/%s", file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer os.Remove(tempPath) // Clean up after processing

	// Add to ChromaDB
	err = services.AddDocuments(c.Request.Context(), "clinic-1", tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PDF successfully added to knowledge base"})
}

func (h *Handler) GetCollection(c *gin.Context) {
	collection := services.GetCollection(c.Request.Context(), "vet_knowledge_base")
	c.JSON(http.StatusOK, gin.H{
		"collection": collection,
	})
}

type CreateCollectionInput struct {
	Name string `json:"name"`
}

func (h *Handler) CreateCollection(c *gin.Context) {
	var input CreateCollectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := services.CreateCollection(c.Request.Context(), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Collection created successfully",
		"collection": res,
	})
}

// ChromaDB Response Struct
type QueryResponse struct {
	Documents [][]string `json:"documents"`
}

func (h *Handler) QueryChromaDB(c *gin.Context) {
	var input services.QueryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := services.QueryChromaDB(c.Request.Context(), input.Query, input.CollectionName, input.NResults)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": response.Documents,
	})
}

func (h *Handler) SearchKnowledgeBase(c *gin.Context) {
	var input services.QueryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := services.SearchKnowledgeBase(c.Request.Context(), input.CollectionName, input.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"documents": response,
	})
}
