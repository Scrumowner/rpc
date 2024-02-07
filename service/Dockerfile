# Блок сборки приложения
FROM golang:1.21.5-alpine3.18 as builder

WORKDIR /usr/local/go/src/service

COPY . .

# Сборка исполняемого файла
RUN go mod tidy
RUN go build -o main

# Вторая стадия сборки
FROM alpine:latest

# Копирование исполняемого файла из первой стадии
COPY --from=builder /usr/local/go/src/service/main /main

EXPOSE 1234

CMD ["/main"]