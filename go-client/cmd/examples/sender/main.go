package main

import (
	"context"
	"fmt"
	"log"
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

	// 요청 큐 선언
	if err := client.DeclareQueue("request_queue"); err != nil {
		log.Fatal(err)
	}

	// 요청 큐 바인딩
	if err := client.BindQueue("request_queue", "test_exchange", "request_key"); err != nil {
		log.Fatal(err)
	}

	// 응답 큐 선언
	replyQueue, err := client.DeclareReplyQueue()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("응답 큐 생성됨: %s", replyQueue)

	// 요청 메시지
	message := []byte("Hello from Go sender!")
	log.Printf("Sending request: %s", message)

	// 상관 ID 생성
	correlationID := fmt.Sprintf("go-req-%d", time.Now().UnixNano())

	// 컨텍스트 생성
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 요청 전송
	err = client.PublishWithOptions(ctx, "test_exchange", "request_key", message, rabbitmq.PublishOptions{
		ReplyTo:       replyQueue,
		CorrelationID: correlationID,
	})
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
	log.Println("Request sent successfully!")

	// 응답 대기
	log.Printf("Waiting for response with correlationID: %s", correlationID)
	response, err := client.WaitForResponse(ctx, replyQueue, correlationID)
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}

	log.Printf("Received response: %s", string(response))

	// 다른 방법으로 테스트
	log.Println("\nTesting another request-response pattern:")

	// 두 번째 요청 메시지
	message2 := []byte("Second request from Go!")
	correlationID2 := fmt.Sprintf("go-req-%d", time.Now().UnixNano())

	// 응답 수신 핸들러
	responseHandler := func(msg []byte) error {
		log.Printf("Received response via handler: %s", string(msg))
		return nil
	}

	// 응답 큐 구독
	if err := client.Consume(replyQueue, responseHandler); err != nil {
		log.Fatalf("Failed to set up response consumer: %v", err)
	}

	// 요청 전송
	err = client.PublishWithOptions(ctx, "test_exchange", "request_key", message2, rabbitmq.PublishOptions{
		ReplyTo:       replyQueue,
		CorrelationID: correlationID2,
	})
	if err != nil {
		log.Fatalf("Failed to publish second message: %v", err)
	}
	log.Printf("Second request sent with correlationID: %s", correlationID2)

	// 응답을 받을 시간을 줌
	time.Sleep(2 * time.Second)
}
