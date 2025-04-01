# RabbitMQ Bridge

Go와 Python 간의 메시지 통신을 위한 RabbitMQ 브릿지 프로젝트입니다.

## 프로젝트 구조

```
rabbitmq-bridge/
├── go-client/          # Go RabbitMQ 클라이언트
│   ├── cmd/           # 실행 가능한 예제
│   │   └── examples/
│   │       ├── sender/    # Go 발신자 예제
│   │       └── receiver/  # Go 수신자 예제
│   ├── internal/      # 내부 패키지
│   └── pkg/           # 공개 패키지
│       └── rabbitmq/   # RabbitMQ 클라이언트 라이브러리
│           ├── client.go      # 기본 클라이언트 구현
│           ├── consumer.go    # 메시지 구독 기능
│           └── publisher.go   # 메시지 발행 기능
├── python-client/      # Python RabbitMQ 클라이언트
│   ├── cmd/           # 실행 가능한 예제
│   │   └── examples/
│   │       ├── basic_usage.py  # 기본 사용법 예제
│   │       ├── sender.py       # Python 발신자 예제
│   │       └── receiver.py     # Python 수신자 예제
│   ├── internal/      # 내부 패키지
│   └── pkg/           # 공개 패키지
│       └── rabbitmq/   # RabbitMQ 클라이언트 라이브러리
│           ├── __init__.py   # 패키지 초기화
│           └── client.py     # 클라이언트 구현
```

## 설치 방법

### Go 클라이언트

```bash
cd go-client
go mod tidy
```

### Python 클라이언트

```bash
cd python-client
pip install -r requirements.txt
# 또는 개발 모드로 설치
pip install -e .
```

## 통신 패턴

### 1. 기본 발행-구독 패턴

발신자가 메시지를 발행하고, 수신자가 메시지를 구독하는 단방향 통신 패턴입니다.

```
[Sender] ----> [RabbitMQ] ----> [Receiver]
```

#### Go 발신자, Python 수신자 예제 실행:

```bash
# 터미널 1 (Python 수신자)
cd python-client
python cmd/examples/receiver.py

# 터미널 2 (Go 발신자)
cd go-client
go run cmd/examples/sender/main.go
```

### 2. 요청-응답 패턴

발신자가 요청을 보내고 수신자가 요청을 처리한 후 응답을 보내는 양방향 통신 패턴입니다.

```
[Sender]                           [Receiver]
   |                                   |
   |-- publish(msg) -----------------> |
   |     replyTo: response_queue       |
   |     correlationId: abc-123        |
   |                                   |
   | <--- publish(response) ---------- |
   |      to: response_queue           |
   |      correlationId: abc-123       |
```

#### Python 수신자, Go 발신자 예제 실행:

```bash
# 터미널 1 (Python 수신자)
cd python-client
python cmd/examples/receiver.py

# 터미널 2 (Go 발신자)
cd go-client
go run cmd/examples/sender/main.go
```

## Go 클라이언트 주요 기능

### 기본 클라이언트

- `NewClient`: RabbitMQ 클라이언트 생성
- `Close`: 연결 종료
- `DeclareExchange`: Exchange 선언
- `DeclareQueue`: Queue 선언
- `BindQueue`: Queue를 Exchange에 바인딩

### 메시지 발행

- `Publish`: 기본 메시지 발행
- `PublishWithOptions`: 옵션(replyTo, correlationId)을 포함한 메시지 발행
- `PublishResponse`: 응답 메시지 발행
- `DeclareReplyQueue`: 임시 응답 큐 선언

### 메시지 구독

- `Consume`: 기본 메시지 구독
- `ConsumeWithContext`: 컨텍스트 기반 메시지 구독
- `ConsumeRequests`: 요청 메시지 구독 및 응답 자동 발행
- `WaitForResponse`: 특정 상관 ID에 대한 응답 대기

## Python 클라이언트 주요 기능

### 기본 클라이언트

- `RabbitMQClient`: RabbitMQ 클라이언트 클래스
- `close`: 연결 종료
- `declare_exchange`: Exchange 선언
- `declare_queue`: Queue 선언
- `bind_queue`: Queue를 Exchange에 바인딩

### 메시지 발행

- `publish`: 메시지 발행 (옵션: reply_to, correlation_id)
- `publish_response`: 응답 메시지 발행
- `declare_reply_queue`: 임시 응답 큐 선언

### 메시지 구독

- `consume`: 기본 메시지 구독
- `consume_with_properties`: 메시지 속성과 함께 구독
- `consume_requests`: 요청 메시지 구독 및 응답 자동 발행
- `rpc_call`: RPC 스타일 요청-응답 수행

## 요구사항

- Go 1.16 이상
- Python 3.7 이상
- RabbitMQ 서버

### RabbitMQ 서버 실행 (Docker)

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

## 라이선스

MIT License
