package messaging

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func InitRabbitMQ() *RabbitClient {
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}
	
	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "guest"
	}
	
	password := os.Getenv("RABBITMQ_PASSWORD")
	if password == "" {
		password = "guest"
	}

	uri := fmt.Sprintf("amqp://%s:%s@%s:5672/", user, password, host)

	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatalf("❌ Impossible de se connecter à RabbitMQ : %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatalf("❌ Impossible d'ouvrir un channel RabbitMQ : %v", err)
	}

	return &RabbitClient{Conn: conn, Channel: ch}
}

func (r *RabbitClient) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}
