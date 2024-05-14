FROM golang:1.22.1 as gobuild

WORKDIR /app

COPY . .

RUN go mod download

ENV CGO_EBABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o ./TODO ./cmd/


FROM alpine

WORKDIR /app

COPY --from=gobuild /app/TODO ./TODO
COPY --from=gobuild /app/web ./web

EXPOSE 7540

ENV TODO_DBFILE=scheduler.db TODO_PASSWORD=@WSX2wsx TODO_PORT=7540

ENTRYPOINT ["./TODO"]


#docker build -t todo_app:v1 .
#docker run -P --rm todo_app:v1