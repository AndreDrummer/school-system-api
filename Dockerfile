FROM golang:1.23.2

WORKDIR /app

# Copy the Go module files (go.mod and go.sum) first to cache dependencies
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o app ./cmd

EXPOSE 8383

CMD [ "./app" ]