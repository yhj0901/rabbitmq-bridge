package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config RabbitMQ 연결 설정을 위한 구조체
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	VHost    string
}

// Client RabbitMQ 클라이언트 구조체
type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewClient 새로운 RabbitMQ 클라이언트를 생성합니다.
func NewClient(config Config) (*Client, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Client{
		conn: conn,
		ch:   ch,
	}, nil
}

// Close 연결을 종료합니다.
func (c *Client) Close() error {
	if err := c.ch.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}
