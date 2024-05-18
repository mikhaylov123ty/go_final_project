#Собираем приложение в образе с golang
FROM golang:1.22.1 as gobuild

WORKDIR /app

COPY . .

RUN go mod download

ENV CGO_EBABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o ./TODO ./cmd/

#Передаем собраный билд в "легкий" и с меньшими уязвимостями alpine и собираем в нем
FROM alpine

WORKDIR /app

COPY --from=gobuild /app/TODO ./TODO
COPY --from=gobuild /app/web ./web

EXPOSE 7540

#Не безопасно, в идеале передал бы через .env файл на этапе запуска контейнера
#Но в ТЗ посоветовали, что можно указать переменные окружения тут
ENV TODO_DBFILE=scheduler.db TODO_PASSWORD=@WSX2wsx TODO_PORT=7540

ENTRYPOINT ["./TODO"]

#Билд образа и запуск контейнера
#docker build -t todo_app:v1 .
#docker run -P --rm todo_app:v1