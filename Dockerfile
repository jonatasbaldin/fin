FROM golang:1.11-alpine
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -ldflags '-extldflags "-static"' -o fin

ENV PATH "$PATH:/app"
CMD ["sh", "-c", "fin -migrate && fin -serve"]