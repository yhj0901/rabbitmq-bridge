# RabbitMQ Bridge Client for Python

Python에서 RabbitMQ를 사용하여 메시지 통신을 쉽게 구현할 수 있는 클라이언트 라이브러리입니다.

## 설치 방법

```bash
# PyPI에서 설치
pip install rabbitmq-bridge-client

# 또는 GitHub에서 직접 설치
pip install git+https://github.com/yhj0901/rabbitmq-bridge.git#subdirectory=python-client
```

## 기본 사용법

### 메시지 발행

```python
from pkg.rabbitmq import RabbitMQClient, Config

# 클라이언트 설정
config = Config(
    host="localhost",
    port=5672,
    username="guest",
    password="guest",
    vhost="/"
)

# 클라이언트 생성 및 사용
with RabbitMQClient(config) as client:
    # Exchange와 Queue 선언
    client.declare_exchange("test_exchange", "direct")
    client.declare_queue("test_queue")
    client.bind_queue("test_queue", "test_exchange", "test_key")
    
    # 메시지 발행
    client.publish("test_exchange", "test_key", "Hello, RabbitMQ!")
```

### 메시지 구독

```python
from pkg.rabbitmq import RabbitMQClient, Config

def message_handler(message):
    print(f"Received message: {message}")

# 클라이언트 설정
config = Config(
    host="localhost",
    port=5672,
    username="guest",
    password="guest",
    vhost="/"
)

# 클라이언트 생성 및 사용
with RabbitMQClient(config) as client:
    # Exchange와 Queue 선언
    client.declare_exchange("test_exchange", "direct")
    client.declare_queue("test_queue")
    client.bind_queue("test_queue", "test_exchange", "test_key")
    
    # 메시지 구독
    client.consume("test_queue", message_handler)
```

## 요청-응답 패턴

### 요청자 (Client)

```python
from pkg.rabbitmq import RabbitMQClient, Config

# 클라이언트 설정
config = Config(
    host="localhost",
    port=5672,
    username="guest",
    password="guest",
    vhost="/"
)

# 클라이언트 생성 및 사용
with RabbitMQClient(config) as client:
    # Exchange와 Queue 선언
    client.declare_exchange("test_exchange", "direct")
    client.declare_queue("request_queue")
    client.bind_queue("request_queue", "test_exchange", "request_key")
    
    # RPC 요청 수행 및 응답 대기
    response = client.rpc_call(
        exchange="test_exchange", 
        routing_key="request_key", 
        message="Hello, RPC request!",
        timeout=5  # 5초 타임아웃
    )
    
    if response:
        print(f"Received response: {response}")
    else:
        print("No response received in time")
```

### 응답자 (Server)

```python
from pkg.rabbitmq import RabbitMQClient, Config

def request_handler(message, reply_to, correlation_id):
    print(f"Received request: {message}")
    # 요청 처리 로직
    return f"Response to: {message}"

# 클라이언트 설정
config = Config(
    host="localhost",
    port=5672,
    username="guest",
    password="guest",
    vhost="/"
)

# 클라이언트 생성 및 사용
with RabbitMQClient(config) as client:
    # Exchange와 Queue 선언
    client.declare_exchange("test_exchange", "direct")
    client.declare_queue("request_queue")
    client.bind_queue("request_queue", "test_exchange", "request_key")
    
    # 요청 구독 및 자동 응답
    client.consume_requests("request_queue", request_handler)
```

## 추가 기능

### 응답 큐 생성

```python
# 응답 큐 생성
reply_queue = client.declare_reply_queue()
```

### 속성과 함께 메시지 발행

```python
# 속성과 함께 메시지 발행
client.publish(
    exchange="test_exchange",
    routing_key="test_key",
    message="Message with properties",
    reply_to="response_queue",
    correlation_id="123456789"
)
```

### 속성과 함께 메시지 구독

```python
def handler_with_properties(message, reply_to, correlation_id):
    print(f"Message: {message}")
    print(f"ReplyTo: {reply_to}")
    print(f"CorrelationID: {correlation_id}")

client.consume_with_properties("test_queue", handler_with_properties)
```

## 라이선스

MIT License 