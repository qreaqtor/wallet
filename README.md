# service-wallet

## Конфигурация

Для работы приложения необходимо наличие `config.env` файла в корне репозитория.
Обязательными являются следующие:
```
DATABASE_USER=
DATABASE_PASSWORD=
DATABASE_NAME=
DATABASE_ADDRESS=
```

## Запуск

Запуск бд в контейнере и перезапись `DATABASE_ADDRESS=postgres:5432`

```
make run-dev
```

## Локальный запуск тестов

Запуск бд в контейнере и перезапись `DATABASE_ADDRESS=localhost:5432`

```
make run-postgres
```

Запуск самого приложения

```
go run cmd/service/main.go
```

Запуск тестов

```
go test -v -tags integration ./...
```

## Swagger

После старта сервиса Swagger будет доступен по адресу

```
https://petstore.swagger.io/?url=http://localhost:8080/api/docs
```
