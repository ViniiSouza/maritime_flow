package minion

import (
	"fmt"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func bindTowersQueue() (<-chan amqp.Delivery, error) {
	q, err := config.Configuration.GetRabbitMQChannel().QueueDeclare(
		config.Configuration.GetTowersQueue(),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare the queue: %w", err)
	}

	msgs, err := config.Configuration.GetRabbitMQChannel().Consume(
		q.Name,
		config.Configuration.GetIdAsString(),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register the consumer: %w", err)
	}

	return msgs, nil
}

func bindAuditQueue() error {
	if err := config.Configuration.GetRabbitMQChannel().ExchangeDeclare(
		"requests",
		amqp.ExchangeDirect,
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange \"requests\": %w", err)
	}

	queue, err := config.Configuration.GetRabbitMQChannel().QueueDeclare(
		config.Configuration.GetAuditQueue(),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue \"%s\": %w", config.Configuration.GetAuditQueue(), err)
	}

	config.Configuration.GetRabbitMQChannel().QueueBind(
		queue.Name,
		"requests",
		"requests",
		false,
		nil,
	)

	return nil
}
