# Используем официальный образ Go как базовый
FROM golang:1.21.5-alpine3.18 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app
# Копируем исходники приложения в рабочую директорию
COPY . .
COPY ./controller /usr/local/go/src/hugoproxy-main/proxy/controller
COPY ./gen /usr/local/go/src/hugoproxy-main/proxy/gen
COPY ./logic /usr/local/go/src/hugoproxy-main/proxy/logic
COPY ./middleware /usr/local/go/src/hugoproxy-main/proxy/middleware
COPY ./provider /usr/local/go/src/hugoproxy-main/proxy/provider
COPY ./responder /usr/local/go/src/hugoproxy-main/proxy/responder
COPY ./service /usr/local/go/src/hugoproxy-main/proxy/service
COPY ./static /usr/local/go/src/hugoproxy-main/proxy/static
COPY ./storage /usr/local/go/src/hugoproxy-main/proxy/storage
COPY ./tools /usr/local/go/src/hugoproxy-main/proxy/tools
COPY ./worker-content /usr/local/go/src/hugoproxy-main/proxy/worker-content
COPY ./controllers /usr/local/go/src/hugoproxy-main/proxy/controllers
COPY ./adapter /usr/local/go/src/hugoproxy-main/proxy/adapter
COPY ./migrator /usr/local/go/src/hugoproxy-main/proxy/migrator
# Скачиваем все зависимости
RUN go mod init test
RUN go mod tidy
RUN go build -o main
# Собираем приложение

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/main /main

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]


