
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


3. слайс vs массив, структура слайса, что проиходит при append
4. map, бакеты, миграции, что проиходит при коллизиях
5. пример собеса https://youtu.be/W_ctQFFnzK0?si=oNIBdNI2ARi-SNWM
6. методы синхронизации https://medium.com/german-gorelkin/synchronization-primitives-go-8857747d9660
7. Data Race vs Race condition

9. Отиличие виртуализации от контейнеризации

10. горутны, сколько вести от и до, GMP, scheduler
11. interface, solid, как в golang опп релизуется

12. Распилить на папки и по архитектуре

13. stack vs heap (go in memory статейка вроде)
14. Индексы, join, eplane/explane analise, окконные фунции, b-tree