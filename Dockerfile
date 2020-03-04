FROM golang:1.13.5

WORKDIR /app

COPY src/go.mod .
COPY src/go.sum .

RUN go mod download

COPY src .

RUN go install main.go

CMD main 
