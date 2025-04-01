import time
from rabbitmq_client.client import RabbitMQClient, Config

def message_handler(message: str):
    print(f"Received message: {message}")

def main():
    # RabbitMQ 클라이언트 설정
    config = Config(
        host="localhost",
        port=5672,
        username="guest",
        password="guest",
        vhost="/"
    )

    # 클라이언트 생성 및 컨텍스트 매니저 사용
    with RabbitMQClient(config) as client:
        # Exchange 선언
        client.declare_exchange("test_exchange", "direct")

        # Queue 선언
        client.declare_queue("test_queue")

        # Queue를 Exchange에 바인딩
        client.bind_queue("test_queue", "test_exchange", "test_key")

        # 메시지 구독 시작 (별도 스레드에서 실행)
        import threading
        consumer_thread = threading.Thread(
            target=client.consume,
            args=("test_queue", message_handler)
        )
        consumer_thread.daemon = True
        consumer_thread.start()

        # 메시지 발행
        message = "Hello, RabbitMQ!"
        client.publish("test_exchange", "test_key", message)

        # 프로그램이 종료되지 않도록 대기
        time.sleep(2)

if __name__ == "__main__":
    main() 