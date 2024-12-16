FROM golang:1.22

WORKDIR /go/src/app
COPY cmd/ cmd/
COPY internal/ internal/
COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o /go/bin/app cmd/main.go

CMD ["/go/bin/app"]
