
Запилить DELETE aviation делай его сам по аналогии с другими ручками (без гпт пж)

Логировать ошибки log.Error(err) чтоб мы понимали что пошло не так (высвечивается в терминале)

Сделать так чтобы сервис запускался в контейнере
В docker-compose.yaml
```golang
api-gateway:
    build: ./
    ports:
      - "8080:8080"  
    depends_on:
      postgres:
        condition: service_started
```
        
добавить ./Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/api-gateway ./
RUN ls -l /app

FROM alpine:latest 

WORKDIR /app
COPY --from=builder /app/api-gateway .

ENTRYPOINT ["./api-gateway"]
```