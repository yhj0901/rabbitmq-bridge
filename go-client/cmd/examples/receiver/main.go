package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yhj0901/rabbitmq-bridge/go-client/pkg/rabbitmq"
)

func main() {
	// 로컬 RabbitMQ 서버 설정
	config := rabbitmq.Config{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		VHost:    "/",
	}

	client, err := rabbitmq.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Exchange 선언
	if err := client.DeclareExchange("test_exchange", "direct"); err != nil {
		log.Fatal(err)
	}

	// Queue 선언
	if err := client.DeclareQueue("request_queue"); err != nil {
		log.Fatal(err)
	}

	// Queue 바인딩
	if err := client.BindQueue("request_queue", "test_exchange", "request_key"); err != nil {
		log.Fatal(err)
	}

	// 컨텍스트 생성
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 요청 처리 핸들러
	requestHandler := func(message []byte, replyTo string, correlationID string) ([]byte, error) {
		log.Printf("Received request: %s", string(message))
		log.Printf("ReplyTo: %s, CorrelationID: %s", replyTo, correlationID)

		// 요청 처리 로직
		time.Sleep(500 * time.Millisecond) // 처리 시간 시뮬레이션

		// 응답 메시지 생성
		response := []byte("Response from Go receiver: " + string(message))
		return response, nil
	}

	// 요청 구독 시작
	if err := client.ConsumeRequests(ctx, "request_queue", requestHandler); err != nil {
		log.Fatal(err)
	}

	log.Println("Waiting for requests. To exit press CTRL+C")

	// 종료 시그널 대기
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
