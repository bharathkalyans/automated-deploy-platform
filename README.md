# GoLang Assignment 3



#### Start ZooKeeper and Kafka Cluster
- `make kafka_up` to start the zookeeper and kafka cluster
- `make kafka_down` to stop the zookeeper and kafka cluster


#### How to start the services
- Open the 2 services in different terminals and make use of Makefile to run the application.

##### Note
- We are using PostgreSQL and Redis (local instances).
- Make sure you are running these 2 before starting the application.
- Also use the following SQL Query to create the Table.
```SQL
    CREATE TABLE IF NOT EXISTS data (
		build_id VARCHAR(255) PRIMARY KEY,
		project_github_url VARCHAR(255) NOT NULL,
		events JSONB NOT NULL
	  );
```