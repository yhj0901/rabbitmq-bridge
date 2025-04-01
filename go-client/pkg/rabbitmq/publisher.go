package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishOptions 메시지 발행 옵션
type PublishOptions struct {
	ReplyTo       string // 응답을 받을 큐 이름
	CorrelationID string // 요청과 응답을 연결하는 상관 ID
}

// Publish 메시지를 발행합니다.
func (c *Client) Publish(ctx context.Context, exchange, routingKey string, message []byte) error {
	return c.PublishWithOptions(ctx, exchange, routingKey, message, PublishOptions{})
}

// PublishWithOptions 옵션과 함께 메시지를 발행합니다.
func (c *Client) PublishWithOptions(ctx context.Context, exchange, routingKey string, message []byte, options PublishOptions) error {
	return c.ch.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          message,
			ReplyTo:       options.ReplyTo,
			CorrelationId: options.CorrelationID,
		})
}

// PublishResponse 응답 메시지를 발행합니다.
func (c *Client) PublishResponse(ctx context.Context, replyTo, correlationID string, message []byte) error {
	return c.ch.PublishWithContext(ctx,
		"",      // default exchange
		replyTo, // 응답 큐 이름
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          message,
			CorrelationId: correlationID,
		})
}

// DeclareExchange exchange를 선언합니다.
func (c *Client) DeclareExchange(name, kind string) error {
	return c.ch.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareQueue queue를 선언합니다.
func (c *Client) DeclareQueue(name string) error {
	_, err := c.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// DeclareReplyQueue 응답 큐를 선언합니다.
func (c *Client) DeclareReplyQueue() (string, error) {
	queue, err := c.ch.QueueDeclare(
		"",    // 이름 자동 생성
		false, // non-durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return "", err
	}
	return queue.Name, nil
}

// BindQueue queue를 exchange에 바인딩합니다.
func (c *Client) BindQueue(queue, exchange, routingKey string) error {
	return c.ch.QueueBind(
		queue,      // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
}
