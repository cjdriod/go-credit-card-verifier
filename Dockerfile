FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-w -extldflags "-static"' -o main ./cmd

## HOSTING server application
#FROM alpine:latest
#
## Create new user and assign /app as home dir
## RUN adduser -S -D -H -h /app appuser
#RUN addgroup -g 1001 appuser && adduser -u 1001 -G appuser -D -h /app appuser
#
#WORKDIR /app
#
#COPY --from=builder /app/main .
#
#RUN chown appuser:appuser main
#RUN chmod +x main
#
## Switch to the app user account
#USER appuser
#
#EXPOSE 8080
#
#CMD ["/app/main"]

# HOSTING server application with empty container
FROM scratch

COPY --from=builder /app/main /main

ENTRYPOINT ["/main"]

