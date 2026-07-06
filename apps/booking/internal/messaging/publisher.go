package messaging

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	PublishBookingCreated(bookingID int, email string) error
}

type rabbitEventPublisher struct {
	channel *amqp.Channel
}

func NewRabbitEventPublisher(ch *amqp.Channel) EventPublisher {
	return &rabbitEventPublisher{channel: ch}
}

func (p *rabbitEventPublisher) PublishBookingCreated(bookingID int, email string) error {
	q, err := p.channel.QueueDeclare("bookings_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"booking_id": bookingID,
		"email":      email,
		"status":     "CREATED",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.channel.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}