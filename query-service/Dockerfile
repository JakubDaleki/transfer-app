FROM golang:1.19
WORKDIR /shared-dependencies
COPY shared-dependencies/ ./
RUN go mod download
WORKDIR /app
COPY query-service/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /query-service
EXPOSE 8888

CMD ["/query-service"]