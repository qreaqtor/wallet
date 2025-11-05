# service-wallet

## Запуск

```
make run-dev
```

## Локальный запуск тестов

```
make run-postgres
```

```
go run cmd/service/main.go
```

```
go test -v -tags integration ./...
```

## Swagger

После старта сервиса Swagger будет доступен по адресу
```
https://petstore.swagger.io/?url=http://localhost:8080/api/docs
```
