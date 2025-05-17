FROM golang:1.23

# Install librdkafka
RUN apt-get update && apt-get install -y librdkafka-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy everything from the root context
COPY . .

# Build the Go app (entry point is in rating/cmd)
RUN go build -o main ./rating/cmd

ENTRYPOINT ["./main"]
