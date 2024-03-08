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

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n\n", msg.TopicPartition, string(msg.Value))
			processBuildEvent(string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}

	}

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
		buildEvent, err := redisClient.RPop(ctx, "build_queue").Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			fmt.Printf("Failed to dequeue build event from Redis queue: %v\n", err)
			continue
		}
		var buildInfo map[string]interface{}
		err = json.Unmarshal([]byte(buildEvent), &buildInfo)
		if err != nil {
			fmt.Printf("Failed to unmarshal build event JSON: %v\n", err)
			return
		}

		// Create an "events" map if it doesn't exist
		events, ok := buildInfo["events"].(map[string]interface{})
		if !ok {
			events = make(map[string]interface{})
			buildInfo["events"] = events
		}

		events["BUILD_QUEUED"] = map[string]interface{}{
			"timestamp": time.Now(),
		}

		fmt.Println("Github URL :: ", buildInfo["project_github_url"])
		concurrentBuilds <- struct{}{}

		// Process build event concurrently
		go func(buildEvent string) {
			defer func() {
				<-concurrentBuilds
			}()
			events["BUILD_STARTED"] = map[string]interface{}{
				"timestamp": time.Now(),
			}
			// Process the build event here
			fmt.Printf("Processing build event: %s\n", buildEvent)

			buildSuccess, portNumber := dockerImplementation(&buildInfo)
			fmt.Printf("Build Success :: %v on Port Number %v", buildSuccess, portNumber)
			if buildSuccess {
				fmt.Println("Build Succeeded !!!!!!")
				events["BUILD_PASSED"] = map[string]interface{}{
					"timestamp": time.Now(),
				}
				events["DEPLOY_PASSED"] = map[string]interface{}{
					"timestamp":          time.Now(),
					"branded_access_url": `https://localhost:8080/` + buildInfo["build_id"].(string),
					"url":                `https://localhost:` + portNumber,
				}
			} else {
				fmt.Println("Build failed !!!!!!")
				events["BUILD_FAILED"] = map[string]interface{}{
					"timestamp": time.Now(),
					"reason":    "GITHUB URL NOT FOUND",
				}
				events["DEPLOY_FAILED"] = map[string]interface{}{
					"timestamp": time.Now(),
					"reason":    "Build Failed",
				}
			}
			saveToPostgres(&buildInfo)

		}(buildEvent)
	}
}

func dockerImplementation(buildInfo *map[string]interface{}) (bool, string) {
	fmt.Println("Inside the Docker Function")
	fmt.Println("build info :: ", (*buildInfo))

	// Extract GitHub URL and build command from buildInfo
	// githubURL, ok := (*buildInfo)["project_github_url"].(string)
	// if !ok {
	// 	fmt.Println("GitHub URL not found in buildInfo")
	// 	return false, ""
	// }
	// buildCommand, ok := (*buildInfo)["build_command"].(string)
	// if !ok {
	// 	fmt.Println("Build command not found in buildInfo")
	// 	return false, ""
	// }

	return true, "1123"
}

func saveToPostgres(buildInfo *map[string]interface{}) {
	buildID, ok := (*buildInfo)["build_id"].(string)
	if !ok {
		fmt.Println("Build ID not found in build event.")
		return
	}

	projectGitHubURL, ok := (*buildInfo)["project_github_url"].(string)
	if !ok {
		fmt.Println("Project GitHub URL not found in build event.")
		return
	}
	eventsJSON, err := json.Marshal((*buildInfo)["events"])

	// fmt.Println("build info :: ", (*buildInfo)["events"])
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
