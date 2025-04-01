import sys
import os
import time

# 프로젝트 루트 디렉토리를 Python 경로에 추가
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../")))

from pkg.rabbitmq import RabbitMQClient, Config

def main():
    # 로컬 RabbitMQ 서버 설정
    config = Config(
        host="localhost",
        port=5672,
        username="guest",
        password="guest",
        vhost="/"
    )

    with RabbitMQClient(config) as client:
        # Exchange 선언
        client.declare_exchange("test_exchange", "direct")

        # 요청 큐 선언
        client.declare_queue("request_queue")

        # 요청 큐 바인딩
        client.bind_queue("request_queue", "test_exchange", "request_key")

        # RPC 요청 수행
        request_message = "Hello from Python!"
        print(f"Sending request: {request_message}")

        # RPC 호출 수행 및 응답 대기
        response = client.rpc_call(
            exchange="test_exchange", 
            routing_key="request_key", 
            message=request_message,
            timeout=5  # 5초 타임아웃
        )

        if response:
            print(f"Received response: {response}")
        else:
            print("No response received in time")

        # 일반 요청-응답 방식으로도 테스트
        print("\nTesting standard request-response pattern:")
        
        # 응답 큐 생성
        reply_queue = client.declare_reply_queue()
        correlation_id = "123456789"

        # 응답 소비자 설정
        def response_handler(message, reply_to, corr_id):
            if corr_id == correlation_id:
                print(f"Received response via callback: {message}")

        # 응답 리스너 시작 (별도 스레드)
        import threading
        consumer_thread = threading.Thread(
            target=client.consume_with_properties,
            args=(reply_queue, response_handler)
        )
        consumer_thread.daemon = True
        consumer_thread.start()

        # 요청 전송
        print(f"Sending request with correlation_id: {correlation_id}")
        client.publish(
            exchange="test_exchange",
            routing_key="request_key",
            message="Manual request from Python",
            reply_to=reply_queue,
            correlation_id=correlation_id
        )

        # 잠시 대기하여 응답을 받을 시간을 줌
        time.sleep(2)

if __name__ == "__main__":
    main() 