package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageHandler 메시지 처리를 위한 함수 타입
type MessageHandler func([]byte) error

// RequestHandler 요청 메시지 처리를 위한 함수 타입
type RequestHandler func([]byte, string, string) ([]byte, error)

// Consume 메시지를 구독합니다.
func (c *Client) Consume(queue string, handler MessageHandler) error {
	msgs, err := c.ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				fmt.Printf("Error handling message: %v\n", err)
			}
			d.Ack(false)
		}
	}()

	return nil
}

// ConsumeWithContext 컨텍스트를 사용하여 메시지를 구독합니다.
func (c *Client) ConsumeWithContext(ctx context.Context, queue string, handler MessageHandler) error {
	msgs, err := c.ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}
				if err := handler(d.Body); err != nil {
					fmt.Printf("Error handling message: %v\n", err)
				}
				d.Ack(false)
			}
		}
	}()

	return nil
}

// ConsumeRequests 요청 메시지를 구독하고 응답을 발행합니다.
func (c *Client) ConsumeRequests(ctx context.Context, queue string, handler RequestHandler) error {
	msgs, err := c.ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}

				// 요청 처리
				response, err := handler(d.Body, d.ReplyTo, d.CorrelationId)
				if err != nil {
					fmt.Printf("Error handling request: %v\n", err)
					d.Ack(false)
					continue
				}

				// 응답이 필요한 경우(ReplyTo가 설정된 경우)에만 응답 발행
				if d.ReplyTo != "" {
					err = c.PublishResponse(ctx, d.ReplyTo, d.CorrelationId, response)
					if err != nil {
						fmt.Printf("Error publishing response: %v\n", err)
					}
				}

				d.Ack(false)
			}
		}
	}()

	return nil
}

// WaitForResponse 응답 큐에서 특정 상관 ID에 해당하는 응답을 기다립니다.
func (c *Client) WaitForResponse(ctx context.Context, replyQueue, correlationID string) ([]byte, error) {
	msgs, err := c.ch.Consume(
		replyQueue, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	responseChannel := make(chan amqp.Delivery)
	go func() {
		for d := range msgs {
			if d.CorrelationId == correlationID {
				responseChannel <- d
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case delivery := <-responseChannel:
		return delivery.Body, nil
	}
}
