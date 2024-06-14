FROM golang:1.15

WORKDIR /app
ADD ./go.mod /app
ADD ./go.sum /app
RUN go mod download

CMD ["sh", "./bin/test.sh"]
