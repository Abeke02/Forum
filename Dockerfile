
# syntax=docker/dockerfile:1
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o cmd/forum cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app .
LABEL version="1.0" 
LABEL creators="@Abdolla && @Bagdaulet02 && @qaharman11 "
EXPOSE 8023
CMD [ "cmd/forum" ]
