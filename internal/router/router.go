package router

import (
	"vet-tails/ai/internal/handlers"
	"vet-tails/ai/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// router.Use(middleware.Authentication())

	// // Services
	// analysisService := services.NewAnalysisService(db)
	soapService := services.NewSOAPService("http://localhost:11434")
	// llavaService := services.NewLlavaService("http://localhost:11434")
	handler := handlers.Handler{
		// DB:          db,
		SoapService: soapService,
		// LlavaService: llavaService,
	}

	// Routes
	api := router.Group("/api/v1")
	{
		// api.POST("/soap", handlers.CreateSOAPNote)
		// api.GET("/patients/:id/history", handlers.GetPatientHistory)
		// api.POST("/invoice/generate", handlers.GenerateInvoice)
		// api.GET("/recommendations", handlers.GetRecommendations)
		api.POST("/soap", handler.CreateSOAPNote)
		// api.POST("/breed", handler.DetectBreed)
		// api.POST("/summary", handler.GeneratePatientSummary)
		// api.POST("/activity", handler.GeneratePetActivityLog)
		api.POST("/upload-pdf", handler.UploadPDFHandler)
		api.GET("/collection", handler.GetCollection)
		api.POST("/collection", handler.CreateCollection)
		api.POST("/search", handler.SearchKnowledgeBase)
	}

	return router
}
