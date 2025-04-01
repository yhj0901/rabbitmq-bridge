import pika
import uuid
import threading
from typing import Callable, Optional, Tuple, Dict
from dataclasses import dataclass
from contextlib import contextmanager

@dataclass
class Config:
    host: str = "localhost"
    port: int = 5672
    username: str = "guest"
    password: str = "guest"
    vhost: str = "/"

class RabbitMQClient:
    def __init__(self, config: Config):
        self.config = config
        self.connection = None
        self.channel = None
        self._connect()
        self.response_callbacks = {}  # correlation_id -> callback

    def _connect(self):
        """RabbitMQ 서버에 연결합니다."""
        credentials = pika.PlainCredentials(
            self.config.username,
            self.config.password
        )
        parameters = pika.ConnectionParameters(
            host=self.config.host,
            port=self.config.port,
            virtual_host=self.config.vhost,
            credentials=credentials
        )
        self.connection = pika.BlockingConnection(parameters)
        self.channel = self.connection.channel()

    def close(self):
        """연결을 종료합니다."""
        if self.channel:
            self.channel.close()
        if self.connection:
            self.connection.close()

    def __enter__(self):
        """with 문에서 사용할 수 있도록 합니다."""
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        """컨텍스트 종료 시 연결을 닫습니다."""
        self.close()

    def declare_exchange(self, name: str, exchange_type: str = "direct"):
        """Exchange를 선언합니다."""
        self.channel.exchange_declare(
            exchange=name,
            exchange_type=exchange_type,
            durable=True
        )

    def declare_queue(self, name: str):
        """Queue를 선언합니다."""
        self.channel.queue_declare(
            queue=name,
            durable=True
        )

    def declare_reply_queue(self) -> str:
        """응답 큐를 선언합니다."""
        result = self.channel.queue_declare(
            queue='',  # 이름 자동 생성
            exclusive=True,
            auto_delete=True
        )
        return result.method.queue

    def bind_queue(self, queue: str, exchange: str, routing_key: str):
        """Queue를 Exchange에 바인딩합니다."""
        self.channel.queue_bind(
            exchange=exchange,
            queue=queue,
            routing_key=routing_key
        )

    def publish(self, exchange: str, routing_key: str, message: str, reply_to: str = None, correlation_id: str = None):
        """메시지를 발행합니다."""
        properties = pika.BasicProperties(
            content_type="text/plain",
            delivery_mode=2,  # 메시지를 영구적으로 저장
            reply_to=reply_to,
            correlation_id=correlation_id
        )
        
        self.channel.basic_publish(
            exchange=exchange,
            routing_key=routing_key,
            body=message,
            properties=properties
        )

    def publish_response(self, reply_to: str, correlation_id: str, message: str):
        """응답 메시지를 발행합니다."""
        properties = pika.BasicProperties(
            content_type="text/plain",
            correlation_id=correlation_id
        )
        
        self.channel.basic_publish(
            exchange='',  # default exchange
            routing_key=reply_to,
            body=message,
            properties=properties
        )

    def consume(self, queue: str, callback: Callable[[str], None]):
        """메시지를 구독합니다."""
        self.channel.basic_consume(
            queue=queue,
            on_message_callback=lambda ch, method, properties, body: callback(body.decode()),
            auto_ack=True
        )
        self.channel.start_consuming()

    def consume_with_properties(self, queue: str, callback: Callable[[str, str, str], None]):
        """속성과 함께 메시지를 구독합니다."""
        def wrapper(ch, method, properties, body):
            reply_to = properties.reply_to or ""
            correlation_id = properties.correlation_id or ""
            callback(body.decode(), reply_to, correlation_id)
            ch.basic_ack(delivery_tag=method.delivery_tag)
            
        self.channel.basic_consume(
            queue=queue,
            on_message_callback=wrapper,
            auto_ack=False
        )
        self.channel.start_consuming()

    def consume_requests(self, queue: str, handler: Callable[[str, str, str], str]):
        """요청 메시지를 구독하고 응답을 발행합니다."""
        def wrapper(ch, method, properties, body):
            # 요청 처리
            message = body.decode()
            reply_to = properties.reply_to
            correlation_id = properties.correlation_id
            
            try:
                # 응답 생성
                response = handler(message, reply_to, correlation_id)
                
                # 응답이 필요한 경우에만 응답 발행
                if reply_to:
                    self.publish_response(reply_to, correlation_id, response)
            except Exception as e:
                print(f"Error handling request: {e}")
            
            ch.basic_ack(delivery_tag=method.delivery_tag)
            
        self.channel.basic_consume(
            queue=queue,
            on_message_callback=wrapper,
            auto_ack=False
        )
        self.channel.start_consuming()

    def setup_response_consumer(self, reply_queue: str):
        """응답 큐 소비자를 설정합니다."""
        def on_response(ch, method, props, body):
            if props.correlation_id in self.response_callbacks:
                callback = self.response_callbacks[props.correlation_id]
                callback(body.decode())
                del self.response_callbacks[props.correlation_id]
        
        self.channel.basic_consume(
            queue=reply_queue,
            on_message_callback=on_response,
            auto_ack=True
        )

    def rpc_call(self, exchange: str, routing_key: str, message: str, timeout: int = 10) -> Optional[str]:
        """RPC 호출을 수행합니다."""
        # 응답 큐 생성
        reply_queue = self.declare_reply_queue()
        
        # 상관 ID 생성
        correlation_id = str(uuid.uuid4())
        
        # 응답을 위한 이벤트와 결과 저장소
        response_event = threading.Event()
        response_result = [None]
        
        # 응답 콜백 등록
        def on_response(response):
            response_result[0] = response
            response_event.set()
        
        self.response_callbacks[correlation_id] = on_response
        
        # 응답 소비자 설정
        self.setup_response_consumer(reply_queue)
        
        # 메시지 발행
        self.publish(
            exchange=exchange,
            routing_key=routing_key,
            message=message,
            reply_to=reply_queue,
            correlation_id=correlation_id
        )
        
        # 응답 대기
        if response_event.wait(timeout):
            return response_result[0]
        else:
            return None

    @contextmanager
    def connection_context(self):
        """컨텍스트 매니저로 연결을 관리합니다."""
        try:
            yield self
        finally:
            self.close() 