FROM golang:1.17-stretch

WORKDIR /src

COPY go.* ./

RUN go mod download

COPY . ./

RUN go mod vendor

RUN go build -v -o main cmd/main.go

CMD [ "./main" ]
