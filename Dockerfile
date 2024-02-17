FROM golang:latest
LABEL authors="Cameron Honis"

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]