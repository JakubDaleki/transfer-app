FROM golang:1.20
WORKDIR /shared-dependencies
COPY shared-dependencies/ ./
RUN go mod download
WORKDIR /app
COPY transfer-aggregator/ ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /transfer-aggregator

CMD ["/transfer-aggregator", "full"]