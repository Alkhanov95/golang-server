
<!-- Запилить DELETE aviation делай его сам по аналогии с другими ручками (без гпт пж) -->

<!-- Логировать ошибки log.Error(err) чтоб мы понимали что пошло не так (высвечивается в терминале) -->

<!-- Сделать так чтобы сервис запускался в контейнере
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
``` -->


1. отличие слайсов от массивов, структура слайса, что происходит при appende
2. мапы, buckets, миграции, колизии





 new
написать ручку по выводу статистики, вывести статистику по id самолета plane_id, вывод: данные по aviation(цена модель) и flights (полеты страна итд)


добавить в папку migrations/ init.sql указать запросы на создании таблицы aviation&flights


сделать так чтобы getplanestatsbyid работал(отдавал в postmane данные по aviation & flights)