# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
# compared to Ubuntu
# FROM golang:1.16-alpine

FROM golang:1.20

WORKDIR /app

# Download necessary Go modules
# COPY go.mod ./
# COPY go.sum ./

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
# COPY ./src .

RUN CGO_ENABLED=0 GOOS=linux go build -o /portimporter
# RUN go build -o /portimporter

# EXPOSE 8080

CMD [ "/portimporter" ]