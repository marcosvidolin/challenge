package main

import (
	"desafio/internal/producer"
	"desafio/internal/service"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		// AMQP connection settings
		amqpURL   = os.Getenv("AMQP_URL")
		amqpQueue = os.Getenv("AMQP_QUEUE")

		// 32 bytes hex key
		cryptorKey = os.Getenv("CRYPTOR_KEY")
	)

	batchSize := flag.Int("b", 100, "Batch size used to send users to the queue")
	file := flag.String("f", "users.csv", "CSV file path")

	flag.Parse()

	reader, err := producer.NewCsvReader(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	cryptor, err := service.NewCryptor(cryptorKey)
	if err != nil {
		log.Fatal(err)
	}

	parser := producer.NewCsvUserParser(cryptor)

	adapter, err := producer.NewRabbitMQAdapter(amqpURL, amqpQueue)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ adapter: %v", err)
	}
	defer adapter.Close()

	up := producer.NewUserProducer(reader, parser, adapter)

	up.Produce(*batchSize)
}
