# Используем официальный образ Go как базовый
FROM golang:1.21.5-alpine3.18 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /usr/local/go/src/proxy
# Копируем исходники приложения в рабочую директорию
COPY . .


# Скачиваем все зависимости
RUN go mod tidy
RUN go build -o main
# Собираем приложение

FROM alpine:latest

# Копирование исполняемого файла из первой стадии
COPY --from=builder /usr/local/go/src/proxy/.env .
COPY --from=builder /usr/local/go/src/proxy/main .
# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]


