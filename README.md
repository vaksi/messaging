# Messaging Service Apps
Messaging service is used by created simple message and retrieve message implementation.

## Cases
1. User Post a massage use API 
2. Message will be store in database 
3. Message `status` should be `sent` if message success in publish to kafka.
4. Message `status` should be `received` if message success to consumed with `consumer_message`.
* Note: I use delay in 5 second in consumer to be update status message `received` for proof.

## Quick Start
### Technology Use
* Database Mysql (MariaDB)
* Kafka (for data streaming)

### Installation Required 
#### Install Golang 
 I use go version for service compatible `go1.12.5 darwin/amd64`
#### Install Docker
I Use Docker for deployment my apps and integration app. Docker is very simple for running this app.

So please install docker into your PC https://docs.docker.com/install

#### Install Kafka and Run
This service using kafka for event streaming. Prepare install please check this link https://kafka.apache.org/quickstart

### Deploy and run 
The firs you must extract this file into your workspace project.

How To run :
    
    Run Mysql
    $ docker-compose -f docker-compose-mysql up   
    
    Run Kafka 
    
    Download dependecies 
    $ go mod download
    
    Run Http Server
    $ go run main.go http
     
    Run Consumer message
    $ go run main.go consumer

## Documentation
This service has documentation following swagger for your try. So please se file in path ./api/api_spech.yaml.

### Unit Test and Mock
To run unit test for the project :

    $ go test $(go list ./... | grep -v /vendor/) -cover
        
### Code Style

In golang they already define the convention style https://golang.org/doc/effective_go.html
