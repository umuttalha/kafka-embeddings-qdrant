package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

type InsertRequest struct {
	Text string `json:"text"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResult struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

type SearchResponse struct {
	Status  string         `json:"status"`
	Results []SearchResult `json:"results"`
}

type Message struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	RequestID string `json:"request_id"`
}

var writer *kafka.Writer
var responseReader *kafka.Reader

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	kafkaHost := os.Getenv("KAFKA_HOST")
	if kafkaHost == "" {
		log.Fatal("KAFKA_HOST environment variable is not set")
	}

	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{fmt.Sprintf("%s:9092", kafkaHost)},
		Topic:   "text-topic",
	})

	responseReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:9092", kafkaHost)},
		Topic:   "response-topic",
		GroupID: "golang-producer",
	})

	// Start response listener
	go listenForResponses()
}

var pendingRequests = make(map[string]chan SearchResponse)

func listenForResponses() {
	for {
		msg, err := responseReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading response: %v", err)
			continue
		}

		var response struct {
			RequestID string         `json:"request_id"`
			Results   SearchResponse `json:"results"`
		}
		if err := json.Unmarshal(msg.Value, &response); err != nil {
			log.Printf("Error unmarshaling response: %v", err)
			continue
		}

		if ch, ok := pendingRequests[response.RequestID]; ok {
			ch <- response.Results
			delete(pendingRequests, response.RequestID)
		}
	}
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	var req InsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := Message{
		Type:    "insert",
		Content: req.Text,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: messageBytes,
		},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Text sent for processing",
	})
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestID := uuid.New().String()
	responseChan := make(chan SearchResponse)
	pendingRequests[requestID] = responseChan

	message := Message{
		Type:      "search",
		Content:   req.Query,
		RequestID: requestID,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: messageBytes,
		},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Wait for response with timeout
	select {
	case response := <-responseChan:
		json.NewEncoder(w).Encode(response)
	case <-time.After(10 * time.Second):
		delete(pendingRequests, requestID)
		http.Error(w, "Search timeout", http.StatusGatewayTimeout)
	}
}

func main() {
	defer writer.Close()

	r := mux.NewRouter()
	r.HandleFunc("/insert", insertHandler).Methods("POST")
	r.HandleFunc("/search", searchHandler).Methods("POST")

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
