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
#### Install Docker
I Use Docker for deployment my apps and integration app. Docker is very simple for running this app.

So please install docker into your PC https://docs.docker.com/install

### Deploy and run 
The firs you must extract this file into your workspace project.

How To run :
    
    $ docker-compose up    

## Documentation
This service has documentation following swagger for your try. So please se file in path ./api/api_spech.yaml.

### Unit Test and Mock
To run unit test for the project :

    $ go test $(go list ./... | grep -v /vendor/) -cover
        
### Code Style

In golang they already define the convention style https://golang.org/doc/effective_go.html
