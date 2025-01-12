FROM golang:1.23.1

WORKDIR /app

COPY . .

RUN chmod +x ./scripts/init_tasks.sh

RUN go build -o main .

CMD ["./main"]