# GoLang Assignment 3


#### Flow of the Application 
https://github.com/bharathkalyans/golang-webknot-3/assets/49526472/90509547-aae7-4be2-974b-f39ce8e03516


#### Start ZooKeeper and Kafka Cluster
- `make kafka_up` to start the zookeeper and kafka cluster
- `make kafka_down` to stop the zookeeper and kafka cluster


##### Note
- We are using PostgreSQL and Redis (local instances).
- Make sure you are running these 2 before starting the application.
- Also use the following SQL Query to create the Table (incase needed, we are calling createTable function to create this table in one of our service).
```SQL
    CREATE TABLE IF NOT EXISTS data (
		build_id VARCHAR(255) PRIMARY KEY,
		project_github_url VARCHAR(255) NOT NULL,
		events JSONB NOT NULL
	  );
```

#### How to start the services
- Open the 2 services in different terminals and make use of Makefile to run the application.
	



#### Postman Documentation
[COLLECTION LINK (JSON)](https://api.postman.com/collections/13437119-a0a7edc9-63f9-47f6-acd5-6e694648b847?access_key=PMAT-01HRF3DDWTFMMCNCTFAGB7FEM6)
