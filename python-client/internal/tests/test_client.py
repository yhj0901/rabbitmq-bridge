import pytest
from pkg.rabbitmq import RabbitMQClient, Config

@pytest.fixture
def config():
    return Config(
        host="localhost",
        port=5672,
        username="guest",
        password="guest",
        vhost="/"
    )

@pytest.fixture
def client(config):
    client = RabbitMQClient(config)
    yield client
    client.close()

def test_declare_exchange(client):
    client.declare_exchange("test_exchange", "direct")
    # Exchange가 선언되었는지 확인
    # 실제로는 exchange_declare가 성공적으로 실행되었는지만 확인

def test_declare_queue(client):
    client.declare_queue("test_queue")
    # Queue가 선언되었는지 확인
    # 실제로는 queue_declare가 성공적으로 실행되었는지만 확인

def test_bind_queue(client):
    client.declare_exchange("test_exchange", "direct")
    client.declare_queue("test_queue")
    client.bind_queue("test_queue", "test_exchange", "test_key")
    # Queue가 Exchange에 바인딩되었는지 확인
    # 실제로는 queue_bind가 성공적으로 실행되었는지만 확인

def test_publish_and_consume(client):
    received_messages = []
    
    def message_handler(message):
        received_messages.append(message)
    
    # Exchange와 Queue 설정
    client.declare_exchange("test_exchange", "direct")
    client.declare_queue("test_queue")
    client.bind_queue("test_queue", "test_exchange", "test_key")
    
    # 메시지 발행
    test_message = "Hello, RabbitMQ!"
    client.publish("test_exchange", "test_key", test_message)
    
    # 메시지 구독
    import threading
    consumer_thread = threading.Thread(
        target=client.consume,
        args=("test_queue", message_handler)
    )
    consumer_thread.daemon = True
    consumer_thread.start()
    
    # 잠시 대기하여 메시지가 처리될 시간을 줌
    import time
    time.sleep(1)
    
    assert test_message in received_messages 