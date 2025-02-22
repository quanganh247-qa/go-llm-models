package main

import (
	"log"
	"vet-tails/ai/internal/router"
)

func main() {
	// Load configuration
	// config := config.LoadConfig()

	// Initialize database
	// db := database.InitDB(config.DatabaseURL)

	// // Initialize services
	// ollamaService := services.NewOllamaService(config.OllamaURL)
	// analysisService := services.NewAnalysisService(db, ollamaService)

	// Setup router
	router := router.SetupRouter()

	// Start server
	log.Fatal(router.Run(":8080"))
}

// package main

// func randomString(n int) string {
// 	b := make([]byte, n)
// 	for i := range b {
// 		b[i] = 'a' + byte(rand.Intn(26))
// 	}
// 	return string(b)
// }

// func main() {
// 	ctx := context.Background()

// 	// 1. Initialize client with proper configuration
// 	client, err := chroma.NewClient(chroma.WithBasePath("http://localhost:8000"))
// 	if err != nil {
// 		fmt.Printf("Failed to create client: %v", err)
// 	}
// 	if err != nil {
// 		log.Fatalf("❌ Failed to create client: %v", err)
// 	}
// 	log.Println("✅ ChromaDB client created successfully")

// 	// 2. Verify server connection
// 	version, err := client.Version(ctx)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to connect to ChromaDB server: %v", err)
// 	}
// 	log.Printf("✅ Connected to ChromaDB server version: %s", version)

// 	// 3. Initialize embedding function with proper error handling
// 	ef, err := ollama.NewOllamaEmbeddingFunction(
// 		ollama.WithBaseURL("http://127.0.0.1:11434"),
// 		ollama.WithModel("nomic-embed-text"),
// 	)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to create embedding function: %v", err)
// 	}
// 	log.Println("✅ Embedding function created successfully")

// 	// embedding function example
// 	embedding, err := ef.EmbedQuery(ctx, "Hello, world!")
// 	if err != nil {
// 		log.Fatalf("❌ Failed to embed text: %v", err)
// 	}
// 	log.Printf("✅ Embedding created successfully: %v", embedding.Len())

// 	// 4. Check if collection exists before creating
// 	collectionName := randomString(10)

// 	// First try to get the collection
// 	existingCollection, err := client.GetCollection(ctx, collectionName, ef)
// 	if err != nil {
// 		// If collection doesn't exist, create it
// 		_, err = client.CreateCollection(ctx, collectionName, nil, false, ef, types.L2)
// 		if err != nil {
// 			log.Printf("❌ Failed to create collection: %v", err)
// 			return
// 		}
// 		log.Printf("✅ Successfully created new collection: %s", collectionName)
// 	} else {
// 		log.Printf("✅ Collection %s already exists", collectionName)
// 	}

// 	if existingCollection != nil {
// 		log.Printf("✅ Collection %s already exists", collectionName)
// 		return
// 	}

// 	// newCollection, err := client.CreateCollection(ctx, collectionName, nil, false, ef, types.L2)
// 	// if err != nil {
// 	// 	log.Printf("❌ Failed to create collection: %v", err)
// 	// 	return
// 	// }
// 	// log.Printf("✅ Successfully created new collection: %s", newCollection.Name)

// }
