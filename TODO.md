
Добавить compose.yaml (docker compose)
- должен запускать postgres в контейнере (почитать что такое docker и контейнеры) 
- когда заработает postgres в контейнере (можно будет подрубится через pgadmin по 0.0.0.0:5432) добавить Dockerfile с сборкой сервиса и запуском в контейнере

### ./compose.yaml
```yaml
services:
  postgres:
    image: postgres:14.10-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_HOST_AUTH_METHOD: trust
    volumes:
      - pgdata:/var/lib/postgresql/data  

volumes:
  pgdata:
```

поставить docker + docker hub