package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yhj0901/rabbitmq-bridge/go-client/pkg/rabbitmq"
)

func main() {
	// RabbitMQ 클라이언트 설정
	config := rabbitmq.Config{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		VHost:    "/",
	}

	// 클라이언트 생성
	client, err := rabbitmq.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer client.Close()

	// Exchange 선언
	if err := client.DeclareExchange("test_exchange", "direct"); err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// Queue 선언
	if err := client.DeclareQueue("test_queue"); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Queue를 Exchange에 바인딩
	if err := client.BindQueue("test_queue", "test_exchange", "test_key"); err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	// 메시지 구독 시작
	handler := func(msg []byte) error {
		fmt.Printf("Received message: %s\n", string(msg))
		return nil
	}

	if err := client.Consume("test_queue", handler); err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}

	// 메시지 발행
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := []byte("Hello, RabbitMQ!")
	if err := client.Publish(ctx, "test_exchange", "test_key", message); err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	// 프로그램이 종료되지 않도록 대기
	time.Sleep(2 * time.Second)
}
