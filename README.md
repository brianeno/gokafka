# go-kafka-example

https://medium.com/swlh/apache-kafka-with-golang-227f9f2eb818


# To Run
## Start Kafka
docker-compose up -d  # remove -d to start in foreground
docker-compose down   # to shut down

## Build Applications
go build -o ./producerapp ./producer/producer.go
go build -o ./consumerapp ./consumer/consumer.go