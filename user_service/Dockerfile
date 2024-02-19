# Первая стадия сборки
FROM golang:1.21.5-alpine3.18 as builder

WORKDIR /usr/local/go/src/user_service

COPY . .

RUN go mod tidy
RUN go build -o main


# Вторая стадия сборки
FROM alpine:latest

# Создайте рабочую директорию
WORKDIR /app

# Скопируйте исполняемый файл из первой стадии сборки
COPY --from=builder /usr/local/go/src/user_service/main .

# Установите права на выполнение для исполняемого файла


EXPOSE 1237

CMD ["./main"]