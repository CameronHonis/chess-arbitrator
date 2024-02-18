FROM golang:latest
LABEL authors="Cameron Honis"

WORKDIR /app
COPY . .
RUN go mod download
RUN go install github.com/onsi/ginkgo/v2/ginkgo

ENV ENV=test
RUN ginkgo -r -v

ENV ENV=prod
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]