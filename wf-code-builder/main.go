// wf-code-builder

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var ctx = context.Background()
var redisClient *redis.Client
var concurrentBuilds = make(chan struct{}, 3)
var dockerClient *client.Client

type BuildEvent struct {
	BuildID          string           `json:"build_id"`
	ProjectGitHubURL string           `json:"project_github_url"`
	Events           map[string]Event `json:"events"`
}

type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason,omitempty"`
	URL       string    `json:"url,omitempty"`
}

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

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server address
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	r := mux.NewRouter()

	fmt.Println("wf-code-builder started at PORT :: 8082")

	go func() {
		log.Fatal(http.ListenAndServe(":8082", r))
	}()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "wf-code-builder-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}

	// Subscribe to the topic
	err = consumer.SubscribeTopics([]string{"wf-code-builder-topic"}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s\n", err)
	}

	go startBuildProcessor()

	// Start consuming messages
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n\n", msg.TopicPartition, string(msg.Value))
			processBuildEvent(string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}

	}
	// consumer.Close()

}

func processBuildEvent(buildEvent string) {
	// Add build event to the Redis queue
	err := redisClient.LPush(ctx, "build_queue", buildEvent).Err()
	if err != nil {
		fmt.Printf("Failed to add build event to Redis queue: %v\n", err)
	}
}

func startBuildProcessor() {
	for {
		// Dequeue build events from Redis
		buildEvent, err := redisClient.RPop(ctx, "build_queue").Result()
		if err == redis.Nil {
			// No build events in the queue, wait or perform other tasks
			continue
		} else if err != nil {
			fmt.Printf("Failed to dequeue build event from Redis queue: %v\n", err)
			continue
		}

		// Acquire a semaphore to control concurrency
		concurrentBuilds <- struct{}{}

		// Process build event concurrently
		go func(buildEvent string) {
			defer func() {
				// Release the semaphore when processing is done
				<-concurrentBuilds
			}()

			// Process the build event here
			fmt.Printf("Processing build event: %s\n", buildEvent)

			// Use Docker API to generate a build with the given build command
			// Implement your Docker API logic here
			fmt.Println("Using Docker API to process it .... .... ...")
			// Simulate build success or error
			containers, err := listContainers()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Containers:")
			for _, container := range containers {
				fmt.Printf("ID: %s, Image: %s, State: %s\n", container.ID, container.Image, container.State)
			}

			buildSuccess := true
			if buildSuccess {
				// Save build information to PostgreSQL
				saveToPostgres(buildEvent)
			} else {
				// Handle build error
				fmt.Println("Build failed !!!!!!")
			}

		}(buildEvent)
	}
}

func saveToPostgres(buildEvent string) {
	var buildInfo map[string]interface{}
	err := json.Unmarshal([]byte(buildEvent), &buildInfo)
	if err != nil {
		fmt.Printf("Failed to unmarshal build event JSON: %v\n", err)
		return
	}

	buildID, ok := buildInfo["build_id"].(string)
	if !ok {
		fmt.Println("Build ID not found in build event.")
		return
	}

	projectGitHubURL, ok := buildInfo["project_github_url"].(string)
	if !ok {
		fmt.Println("Project GitHub URL not found in build event.")
		return
	}

	eventsJSON, err := json.Marshal(buildInfo["events"])
	if err != nil {
		fmt.Printf("Failed to marshal events JSON: %v\n", err)
		return
	}

	// Save build information to PostgreSQL
	query := "INSERT INTO data (build_id, project_github_url, events) VALUES ($1, $2, $3)"
	_, err = db.Exec(query, buildID, projectGitHubURL, eventsJSON)
	if err != nil {
		fmt.Printf("Failed to save build information to PostgreSQL: %v\n", err)
		return
	}

	fmt.Println("Build information saved to PostgreSQL.")
}

func listContainers() ([]types.Container, error) {
	// Get a list of containers
	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	return containers, nil
}
