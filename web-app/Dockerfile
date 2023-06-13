FROM golang:1.19
WORKDIR /shared-dependencies
COPY shared-dependencies/ ./
RUN go mod download
WORKDIR /app
COPY web-app/ ./
RUN go mod download
WORKDIR /app/cmd/service
RUN CGO_ENABLED=0 GOOS=linux go build -o /web-app
EXPOSE 8000

CMD ["/web-app"]