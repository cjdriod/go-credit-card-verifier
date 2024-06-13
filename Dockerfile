FROM golang:latest as builder

WORKDIR /app

# Copy go.mod and go.sum then download the packages 1st before copy everything from the directory,
# can help on caching doker container, so that next time dont have to do it again if no changes
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-w -extldflags "-static"' -o main ./cmd

# HOSTING server application
FROM alpine:latest

# Create new user and assign /app as home dir
# RUN adduser -S -D -H -h /app appuser
RUN addgroup -g 1001 appuser && adduser -u 1001 -G appuser -D -h /app appuser

WORKDIR /app

COPY --from=builder /app/main .

RUN chown appuser:appuser main
RUN chmod +x main

# Switch to the app user account
USER appuser

EXPOSE 8080

CMD ["/app/main"]

