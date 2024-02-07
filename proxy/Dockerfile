# Используем официальный образ Go как базовый
FROM golang:1.21.5-alpine3.18 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /usr/local/go/src/hugoproxy-main/proxy
# Копируем исходники приложения в рабочую директорию
COPY . .


# Скачиваем все зависимости
RUN go mod init test
RUN go mod tidy
RUN go build -o main
# Собираем приложение

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder  /usr/local/go/src/hugoproxy-main/proxy/main /main

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]


