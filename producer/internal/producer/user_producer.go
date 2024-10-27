package producer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type FileReader interface {
	Read(chunkSize int, chunkCh chan<- Chunk) error
	Close() error
}

type Parser interface {
	Parse(r Record) (*User, error)
}

type userProducer struct {
	reader FileReader
	parser Parser

	amqpAdapter AmqpAdapter
}

type AmqpAdapter interface {
	Publish(message, contentType string) error
}

func NewUserProducer(r FileReader, p Parser, amqpAdapter AmqpAdapter) *userProducer {
	return &userProducer{
		reader:      r,
		parser:      p,
		amqpAdapter: amqpAdapter,
	}
}

// TODO: doc
func (u *userProducer) Produce(chunkSize int) error {
	chunkCh := make(chan Chunk)

	go func() {
		err := u.reader.Read(chunkSize, chunkCh)
		if err != nil && err != io.EOF {
			log.Fatal("Error reading chunk:", err)
		}
	}()

	for c := range chunkCh {
		users := make([]User, 0, len(c.Records))
		for _, r := range c.Records {
			user, err := u.parser.Parse(r)
			if err != nil {
				log.Printf("error parsing record on line %v", err)
				continue
			}
			users = append(users, *user)
		}

		msg, err := json.Marshal(users)
		if err != nil {
			return fmt.Errorf("failed to marshal users to JSON: %w", err)
		}

		// TODO: error handler
		u.amqpAdapter.Publish(string(msg), "application/json")
	}

	return nil
}
