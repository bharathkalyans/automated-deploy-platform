package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"collect-client/models"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	BuildQueued  = "BUILD_QUEUED"
	BuildStarted = "BUILD_STARTED"
	BuildFailed  = "BUILD_FAILED"
	BuildPassed  = "BUILD_PASSED"
	DeployPassed = "DEPLOY_PASSED"
	DeployFailed = "DEPLOY_FAILED"
)

var db *sql.DB

func init() {
	connStr := "postgresql://postgres:@localhost:5432/golang?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		log.Fatal("Unable to Connect to PostgresDB ....")
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB connection succesful ....")
	createDataTable(db)
}

func main() {

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/api/v1/collect", CollectHandler).Methods("POST")
	r.HandleFunc("/api/v1/build/{build_id}", BuildInfoHandler).Methods("GET")

	fmt.Println("Server started at PORT :: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func BuildInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	buildID := vars["build_id"]
	fmt.Println("Build ID :: ", buildID)

	query := "SELECT * FROM data WHERE build_id = $1"
	row := db.QueryRow(query, buildID)

	var buildInfo models.BuildInfo
	err := row.Scan(&buildInfo.BuildID, &buildInfo.ProjectGithubURL, &buildInfo.Events)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(buildInfo)
	}
}

func CollectHandler(w http.ResponseWriter, r *http.Request) {
	var collectRequest models.CollectRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&collectRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Request Body %+v", collectRequest)

	defer r.Body.Close()

	buildID := uuid.New().String()

	// Send message to Kafka Cluster
	err = sendToKafka(buildID, collectRequest)
	if err != nil {
		http.Error(w, "Failed to send message to Kafka", http.StatusInternalServerError)
		return
	}

	// Return the build ID to the user
	response := models.CollectResponse{BuildID: buildID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendToKafka(buildID string, request models.CollectRequest) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return err
	}
	defer p.Close()

	topic := "wf-code-builder-topic"

	message := fmt.Sprintf(`{"build_id": "%s", "project_github_url": "%s", "build_command": "%s", "build_out_dir": "%s"}`,
		buildID, request.ProjectGithubURL, request.BuildCommand, request.BuildOutDir)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		return err
	}
	fmt.Println(" \nSent the Topic to Kafka Server .....")
	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
	// p.Close()

	return nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\nRequest: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func createDataTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS data (
		build_id VARCHAR(255) PRIMARY KEY,
		project_github_url VARCHAR(255) NOT NULL,
		events JSONB NOT NULL
	  );	  
	`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

}
