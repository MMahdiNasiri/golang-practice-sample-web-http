FROM golang:1.25-alpine

WORKDIR /app

RUN apk --no-cache add git && \
    go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air"]
