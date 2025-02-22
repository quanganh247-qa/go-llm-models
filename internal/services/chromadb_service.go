package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/ledongthuc/pdf"
)

const (
	ChromaURL = "http://localhost:8000/api/v1"
	// Adjust these constants based on your needs
	MaxChunkSize    = 1500  // Maximum characters per chunk
	MinChunkSize    = 500   // Minimum characters per chunk
	OverlapSize     = 100   // Number of words to overlap between chunks
	SentenceEndings = ".!?" // Characters that denote sentence endings
)

type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

type QueryRequest struct {
	CollectionName string `json:"collection_name"`
	Query          string `json:"query"`
	NResults       int    `json:"n_results"`
}

type Document struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Source string `json:"source"`
}

type AddRequest struct {
	Collection string      `json:"collection"`
	Embeddings [][]float32 `json:"embeddings"`
	Documents  []Document  `json:"documents"`
	IDs        []string    `json:"ids"`
}

func readPDF(filepath string) (string, error) {
	f, r, err := pdf.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("error opening PDF: %v", err)
	}
	defer f.Close()

	var text string
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		content, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		text += content
	}

	return text, nil
}

func ChromaClient() *chroma.Client {
	client, err := chroma.NewClient()
	if err != nil {
		log.Fatalf("Error creating client: %s \n", err)
		return nil
	}
	return client
}

type Collection struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

func GetCollection(ctx context.Context, collectionName string) *Collection {
	client := ChromaClient()

	ef, err := ollama.NewOllamaEmbeddingFunction(ollama.WithBaseURL("http://127.0.0.1:11434"), ollama.WithModel("nomic-embed-text"))
	if err != nil {
		fmt.Printf("Error creating Ollama embedding function: %s \n", err)
	}

	// Try to get existing collection
	collection, err := client.GetCollection(ctx, collectionName, ef)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &Collection{
		ID:       collection.ID,
		Name:     collection.Name,
		Metadata: collection.Metadata,
	}
}

func CreateCollection(ctx context.Context, collectionName string) (*Collection, error) {
	log.Printf("📝 Attempting to create collection: %s", collectionName)

	client := ChromaClient()
	if client == nil {
		log.Println("❌ Failed to create ChromaDB client")
		return nil, fmt.Errorf("failed to create ChromaDB client")
	}
	log.Println("✅ ChromaDB client created successfully")

	ef, err := ollama.NewOllamaEmbeddingFunction(
		ollama.WithBaseURL("http://127.0.0.1:11434"),
		ollama.WithModel("nomic-embed-text"),
	)
	if err != nil {
		log.Printf("❌ Error creating Ollama embedding function: %s\n", err)
		return nil, fmt.Errorf("failed to create Ollama embedding function")
	}
	log.Println("✅ ChromaDB client created successfully")

	collection, err := client.CreateCollection(ctx, collectionName, nil, true, ef, types.L2)
	if err != nil {
		log.Printf("❌ Error creating collection: %s\n", err)
		return nil, fmt.Errorf("failed to create collection: %v", err)
	}
	log.Printf("✅ Successfully created new collection: %s", collectionName)

	return &Collection{
		ID:       collection.ID,
		Name:     collection.Name,
		Metadata: collection.Metadata,
	}, nil
}

func splitContent(content string) []string {
	var chunks []string

	// Normalize content: remove excessive whitespace
	content = strings.Join(strings.Fields(content), " ")

	// Split into sentences first
	sentences := splitIntoSentences(content)

	currentChunk := strings.Builder{}

	for i := 0; i < len(sentences); i++ {
		sentence := sentences[i]

		// If adding this sentence would exceed MaxChunkSize and we're above MinChunkSize
		if currentChunk.Len()+len(sentence) > MaxChunkSize && currentChunk.Len() >= MinChunkSize {
			// Store the current chunk
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))

			// Start new chunk with overlap from previous chunk
			currentChunk.Reset()

			// Add overlap from previous chunk if possible
			if lastChunk := chunks[len(chunks)-1]; len(lastChunk) > 0 {
				words := strings.Fields(lastChunk)
				if len(words) > OverlapSize {
					overlap := strings.Join(words[len(words)-OverlapSize:], " ")
					currentChunk.WriteString(overlap)
					currentChunk.WriteString(" ")
				}
			}
		}

		currentChunk.WriteString(sentence)
		currentChunk.WriteString(" ")
	}

	// Add the final chunk if it's not empty
	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}

// Helper function to split text into sentences
func splitIntoSentences(text string) []string {
	var sentences []string
	var current strings.Builder

	for i := 0; i < len(text); i++ {
		current.WriteByte(text[i])

		// Check if we're at a sentence ending
		if strings.ContainsRune(SentenceEndings, rune(text[i])) {
			// Look ahead to handle ellipsis and other edge cases
			if i+1 < len(text) && strings.ContainsRune(SentenceEndings, rune(text[i+1])) {
				continue
			}

			// Add sentence if it's not empty
			if sentence := strings.TrimSpace(current.String()); len(sentence) > 0 {
				sentences = append(sentences, sentence)
				current.Reset()
			}
		}
	}

	// Add any remaining text as a sentence
	if remaining := strings.TrimSpace(current.String()); len(remaining) > 0 {
		sentences = append(sentences, remaining)
	}

	return sentences
}

func generateEmbeddings(ctx context.Context, content string) (*types.Embedding, error) {
	ef, err := ollama.NewOllamaEmbeddingFunction(
		ollama.WithBaseURL("http://127.0.0.1:11434"),
		ollama.WithModel("nomic-embed-text"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama embedding function")
	}

	embedding, err := ef.EmbedQuery(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to embed content")
	}

	return embedding, nil
}

func AddDocuments(ctx context.Context, collectionName string, filepath string) error {
	// Khởi tạo ChromaDB client
	client := ChromaClient()
	if client == nil {
		log.Println("❌ Failed to create ChromaDB client")
		return fmt.Errorf("failed to create ChromaDB client")
	}
	log.Println("✅ ChromaDB client created successfully")

	// Tạo embedding function với Ollama
	ef, err := ollama.NewOllamaEmbeddingFunction(
		ollama.WithBaseURL("http://127.0.0.1:11434"),
		ollama.WithModel("nomic-embed-text"),
	)
	if err != nil {
		log.Printf("❌ Error creating Ollama embedding function: %s\n", err)
		return fmt.Errorf("failed to create Ollama embedding function")
	}

	// Lấy collection
	collection, err := client.GetCollection(ctx, collectionName, ef)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
		return fmt.Errorf("failed to create collection")
	}

	// Đọc file PDF
	data, err := readPDF(filepath)
	if err != nil {
		log.Fatalf("Error reading PDF: %v\n", err)
		return fmt.Errorf("failed to read PDF")
	}

	// Split content into chunks
	chunks := splitContent(data)

	for i, chunk := range chunks {
		fmt.Printf("chunk: %v\n", chunk)
		embedding, err := generateEmbeddings(ctx, chunk)
		if err != nil {
			log.Fatalf("Error generating embeddings: %v\n", err)
			return fmt.Errorf("failed to generate embeddings")
		}
		fmt.Printf("embedding: %v\n", embedding)
		metadata := map[string]interface{}{
			"source": filepath,
			"page":   i + 1,
		}
		_, err = collection.Add(
			ctx,
			[]*types.Embedding{embedding},
			[]map[string]interface{}{metadata},
			[]string{chunk},
			[]string{fmt.Sprintf("doc_1_chunk_%d", i+1)},
		)
		if err != nil {
			log.Fatalf("Error adding document: %v\n", err)
			return fmt.Errorf("failed to add document")
		}
		fmt.Printf("Added document: %v\n", chunk)
	}
	return nil
}
func SearchKnowledgeBase(ctx context.Context, collectionName string, query string) ([]string, error) {
	client := ChromaClient()
	if client == nil {
		log.Println("❌ Failed to create ChromaDB client")
		return nil, fmt.Errorf("failed to create ChromaDB client")
	}
	log.Println("✅ ChromaDB client created successfully")

	ef, err := ollama.NewOllamaEmbeddingFunction(
		ollama.WithBaseURL("http://127.0.0.1:11434"),
		ollama.WithModel("nomic-embed-text"),
	)
	if err != nil {
		log.Printf("❌ Error creating Ollama embedding function: %s\n", err)
		return nil, fmt.Errorf("failed to create Ollama embedding function")
	}
	// Truy vấn ChromaDB để tìm các tài liệu liên quan
	collection, err := client.GetCollection(ctx, collectionName, ef)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %v", err)
	}

	qr, err := collection.Query(ctx, []string{query}, 5, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error querying ChromaDB: %v", err)
	}

	// Lấy danh sách các đoạn văn bản liên quan
	relevantTexts := make([]string, 0)
	for _, doc := range qr.Documents {
		relevantTexts = append(relevantTexts, doc[0])
	}

	fmt.Printf("relevantTexts: %v\n", relevantTexts)

	return relevantTexts, nil
}

// Function to Query ChromaDB
func QueryChromaDB(ctx context.Context, query string, collectionName string, nResults int) (*QueryResponse, error) {
	client := ChromaClient()
	if client == nil {
		log.Println("❌ Failed to create ChromaDB client")
		return nil, fmt.Errorf("failed to create ChromaDB client")
	}
	log.Println("✅ ChromaDB client created successfully")

	ef, err := ollama.NewOllamaEmbeddingFunction(
		ollama.WithBaseURL("http://127.0.0.1:11434"),
		ollama.WithModel("nomic-embed-text"),
	)
	if err != nil {
		log.Printf("❌ Error creating Ollama embedding function: %s\n", err)
		return nil, fmt.Errorf("failed to create Ollama embedding function")
	}

	collection, err := client.GetCollection(ctx, collectionName, ef)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %v", err)
	}

	// Perform vector search with configurable number of results
	results, err := collection.Query(
		ctx,
		[]string{query},
		int32(nResults),
		nil, // where
		nil, // whereDocument
		nil, // include
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query collection: %v", err)
	}

	return &QueryResponse{
		Documents: results.Documents,
		Distances: results.Distances,
		Metadatas: results.Metadatas,
	}, nil
}

type QueryResponse struct {
	Documents [][]string                 `json:"documents"`
	Distances [][]float32                `json:"distances"`
	Metadatas [][]map[string]interface{} `json:"metadatas"`
	IDs       [][]string                 `json:"ids"`
}
