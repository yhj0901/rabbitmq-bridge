import sys
import os
import signal
import time

# 프로젝트 루트 디렉토리를 Python 경로에 추가
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../")))

from pkg.rabbitmq import RabbitMQClient, Config

def request_handler(message, reply_to, correlation_id):
    """요청을 처리하고 응답을 생성합니다."""
    print(f"Received request: {message}")
    print(f"ReplyTo: {reply_to}, CorrelationID: {correlation_id}")
    
    # 요청 처리 로직
    time.sleep(0.5)  # 처리 시간 시뮬레이션
    
    # 응답 메시지 생성
    response = f"Response from Python receiver: {message}"
    return response

def main():
    # 로컬 RabbitMQ 서버 설정
    config = Config(
        host="localhost",
        port=5672,
        username="guest",
        password="guest",
        vhost="/"
    )

    client = RabbitMQClient(config)

    def signal_handler(signum, frame):
        print("\nShutting down...")
        client.close()
        sys.exit(0)

    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)

    try:
        # Exchange 선언
        client.declare_exchange("test_exchange", "direct")

        # 요청 큐 선언
        client.declare_queue("request_queue")

        # 요청 큐 바인딩
        client.bind_queue("request_queue", "test_exchange", "request_key")

        print("Waiting for requests. To exit press CTRL+C")
        
        # 요청 구독 시작
        client.consume_requests("request_queue", request_handler)
    except KeyboardInterrupt:
        print("\nShutting down...")
    finally:
        client.close()

if __name__ == "__main__":
    main() 