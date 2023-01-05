FROM golang:1.17-stretch

WORKDIR /src

# Cache dependencies
ADD go.mod .
ADD go.sum .

RUN go mod download

ADD . .

RUN go mod tidy

RUN go build -v -o main cmd/main.go

CMD [ "./main" ]
