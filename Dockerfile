FROM golang:alpine AS builder
WORKDIR /app

# Copy go mod and sum files first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source code
COPY . .

RUN go build -o attendance-service main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/attendance-service .
COPY .env .
EXPOSE 3000
CMD ["./attendance-service"]