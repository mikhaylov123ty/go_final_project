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

ENTRYPOINT ["./TODO"]

#EXPOSE 7540
#docker build -t todo_app:v1 .
#docker run -p 7540:7540 --env-file /Users/amikhaylov/Downloads/.env -it todo_app:v1