// wf-code-builder

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
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
}

func main() {

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

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

	// Start consuming messages
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			// Process the message here
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\nRequest: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
