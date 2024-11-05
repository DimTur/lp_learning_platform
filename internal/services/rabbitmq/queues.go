package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DeclareQueue announces the queue
func (c *RMQClient) DeclareQueue(
	name string,
	durable, autoDelete, exclusive, noWait bool,
	args map[string]interface{}) (amqp.Queue, error) {

	tableArgs := amqp.Table(args)

	return c.adminCH.QueueDeclare(
		name,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		tableArgs,
	)
}

// PublishToQueue send message to queue
func (c *RMQClient) PublishToQueue(ctx context.Context, queueName string, body []byte) error {
	return c.publishWithContext(ctx, "", queueName, body)
}
