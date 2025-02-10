FROM golang:alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN apk add --no-cache zip ffmpeg

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app ./cmd/app/main.go

EXPOSE 3333

CMD ["app"]
